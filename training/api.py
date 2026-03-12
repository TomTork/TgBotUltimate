import os
from typing import Dict, Any, List

from fastapi import FastAPI, HTTPException
from pydantic import BaseModel, Field
from starlette.concurrency import run_in_threadpool

from realty_query_ner import (
    ner_predict_entities,
    build_schema_from_text,
    PredEntity,
    ID2LABEL,
)

from transformers import AutoTokenizer, AutoModelForTokenClassification

# Config

from pathlib import Path

BASE_DIR = Path(__file__).resolve().parent  # .../training
USE_ONNX = os.getenv("USE_ONNX", "1") == "1"

MODEL_DIR = os.getenv("MODEL_DIR", str(BASE_DIR / "onnx"))
DEVICE = os.getenv("DEVICE", "cpu")


try:
    if USE_ONNX:
        from optimum.onnxruntime import ORTModelForTokenClassification
except Exception:
    ORTModelForTokenClassification = None



# API schemas

class ParseRequest(BaseModel):
    text: str = Field(..., min_length=1, description="User query text")


class ParseResponse(BaseModel):
    project_name: Any
    building_liter: Any
    floor_min: Any
    floor_max: Any
    rooms_amount_min: Any
    rooms_amount_max: Any
    square_min: Any
    square_max: Any
    cost_min: Any
    cost_max: Any


class HealthResponse(BaseModel):
    status: str
    model_dir: str
    use_onnx: bool


# App + model state
app = FastAPI(title="API", version="1.0.0")

_tokenizer = None
_model = None


def _load_model():
    global _tokenizer, _model

    if not os.path.exists(MODEL_DIR):
        raise RuntimeError(f"MODEL_DIR not found: {MODEL_DIR}")

    _tokenizer = AutoTokenizer.from_pretrained(MODEL_DIR)

    if USE_ONNX:
        if ORTModelForTokenClassification is None:
            raise RuntimeError("USE_ONNX=1 but optimum.onnxruntime is not available")
        _model = ORTModelForTokenClassification.from_pretrained(MODEL_DIR)
    else:
        _model = AutoModelForTokenClassification.from_pretrained(MODEL_DIR)
        _model.to(DEVICE)
        _model.eval()


@app.on_event("startup")
def startup_event():
    _load_model()


def _predict_sync(text: str) -> Dict[str, Any]:
    if _tokenizer is None or _model is None:
        raise RuntimeError("Model is not loaded")

    if not USE_ONNX:
        entities = ner_predict_entities(text, _tokenizer, _model, device=DEVICE)
        return build_schema_from_text(text, entities)

    import torch

    enc = _tokenizer(
        text,
        return_offsets_mapping=True,
        return_tensors="pt",
        truncation=True,
        max_length=192,
    )
    offsets = enc.pop("offset_mapping")[0].tolist()

    with torch.no_grad():
        out = _model(**enc)
        logits = out.logits[0]
        probs = torch.softmax(logits, dim=-1)
        pred_ids = torch.argmax(probs, dim=-1).cpu().tolist()
        pred_scores = torch.max(probs, dim=-1).values.cpu().tolist()

    entities: List[PredEntity] = []
    cur_label = None
    cur_start = None
    cur_end = None
    cur_scores = []

    for pid, (os_, oe_), sc in zip(pred_ids, offsets, pred_scores):
        if os_ == oe_:
            continue
        lab = ID2LABEL[int(pid)]
        if lab == "O":
            if cur_label is not None:
                avg = float(sum(cur_scores) / max(1, len(cur_scores)))
                entities.append(PredEntity(cur_label, cur_start, cur_end, text[cur_start:cur_end], avg))
                cur_label = None
                cur_start = cur_end = None
                cur_scores = []
            continue

        bio, etype = lab.split("-", 1)
        if bio == "B" or (cur_label != etype):
            if cur_label is not None:
                avg = float(sum(cur_scores) / max(1, len(cur_scores)))
                entities.append(PredEntity(cur_label, cur_start, cur_end, text[cur_start:cur_end], avg))
            cur_label = etype
            cur_start = os_
            cur_end = oe_
            cur_scores = [sc]
        else:
            cur_end = oe_
            cur_scores.append(sc)

    if cur_label is not None:
        avg = float(sum(cur_scores) / max(1, len(cur_scores)))
        entities.append(PredEntity(cur_label, cur_start, cur_end, text[cur_start:cur_end], avg))

    return build_schema_from_text(text, entities)


# Routes
@app.get("/health", response_model=HealthResponse)
async def health():
    return HealthResponse(status="ok", model_dir=MODEL_DIR, use_onnx=USE_ONNX)


@app.post("/parse", response_model=ParseResponse)
async def parse(req: ParseRequest):
    text = req.text.strip()
    if not text:
        raise HTTPException(status_code=400, detail="Empty text")

    result = await run_in_threadpool(_predict_sync, text)
    return result


if __name__ == "__main__":
    import uvicorn

    host = os.getenv("NEURO_HOST", "127.0.0.1")
    port = int(os.getenv("NEURO_PORT", "8000"))

    uvicorn.run(
        app,
        host=host,
        port=port,
        log_level="info",
    )

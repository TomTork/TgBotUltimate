from __future__ import annotations

import argparse
import json
import os
import random
import re
from dataclasses import dataclass
from typing import Dict, List, Tuple, Optional, Any

import numpy as np

from datasets import Dataset
from transformers import (
    AutoTokenizer,
    AutoModelForTokenClassification,
    DataCollatorForTokenClassification,
    TrainingArguments,
    Trainer,
    set_seed,
)

# Labels
ENTITY_TYPES = [
    "PROJECT",   # project_name
    "LITER",     # building_liter
    "FLOOR",     # floor
    "ROOMS",     # rooms
    "SQUARE",    # area
    "COST",      # price
]
LABELS = ["O"] + [f"B-{t}" for t in ENTITY_TYPES] + [f"I-{t}" for t in ENTITY_TYPES]
LABEL2ID = {l: i for i, l in enumerate(LABELS)}
ID2LABEL = {i: l for l, i in LABEL2ID.items()}
DEFAULT_TRAIN_JSON = "train_data.json"


# Демо
def demo_training_samples() -> List[Dict[str, Any]]:
    samples = []

    def ent(text: str, substring: str, label: str) -> Dict[str, Any]:
        start = text.lower().index(substring.lower())
        end = start + len(substring)
        return {"start": start, "end": end, "label": label}

    t1 = "ищу трехкомнатную квартиру 70-85 м2 до 12 млн в ЖК Сосны, литер Б, этаж 5-12"
    samples.append({
        "text": t1,
        "entities": [
            ent(t1, "трехкомнатную", "ROOMS"),
            ent(t1, "70-85 м2", "SQUARE"),
            ent(t1, "до 12 млн", "COST"),
            ent(t1, "ЖК Сосны", "PROJECT"),
            ent(t1, "литер Б", "LITER"),
            ent(t1, "этаж 5-12", "FLOOR"),
        ],
    })

    t2 = "нужна 2-комнатная от 55 кв м, не выше 8 этажа, бюджет 9 500 000 руб, проект Северный"
    samples.append({
        "text": t2,
        "entities": [
            ent(t2, "2-комнатная", "ROOMS"),
            ent(t2, "от 55 кв м", "SQUARE"),
            ent(t2, "не выше 8 этажа", "FLOOR"),
            ent(t2, "9 500 000 руб", "COST"),
            ent(t2, "проект Северный", "PROJECT"),
        ],
    })

    t3 = "хочу студию до 30м2, этаж 2-5, до 4.2млн, в жк 'Панорама' корпус А"
    samples.append({
        "text": t3,
        "entities": [
            ent(t3, "студию", "ROOMS"),
            ent(t3, "до 30м2", "SQUARE"),
            ent(t3, "этаж 2-5", "FLOOR"),
            ent(t3, "до 4.2млн", "COST"),
            ent(t3, "жк 'Панорама'", "PROJECT"),
            ent(t3, "корпус А", "LITER"),
        ],
    })

    t4 = "ищу 4к квартиру 90–110 м², от 15 млн, этаж от 10, жк Олимп"
    samples.append({
        "text": t4,
        "entities": [
            ent(t4, "4к", "ROOMS"),
            ent(t4, "90–110 м²", "SQUARE"),
            ent(t4, "от 15 млн", "COST"),
            ent(t4, "этаж от 10", "FLOOR"),
            ent(t4, "жк Олимп", "PROJECT"),
        ],
    })

    return samples


def load_training_samples(path: str) -> List[Dict[str, Any]]:
    with open(path, "r", encoding="utf-8") as f:
        data = json.load(f)

    if not isinstance(data, list):
        raise ValueError("Training JSON must be a list of samples")

    for i, item in enumerate(data):
        if "text" not in item or "entities" not in item:
            raise ValueError(f"Bad sample #{i}: must contain 'text' and 'entities'")
        if not isinstance(item["text"], str):
            raise ValueError(f"Bad sample #{i}: 'text' must be a string")
        if not isinstance(item["entities"], list):
            raise ValueError(f"Bad sample #{i}: 'entities' must be a list")
        for e in item["entities"]:
            if not all(k in e for k in ("start", "end", "label")):
                raise ValueError(f"Bad entity in sample #{i}: {e}")
            if e["label"] not in ENTITY_TYPES:
                raise ValueError(f"Bad entity label in sample #{i}: {e['label']} not in {ENTITY_TYPES}")
            if not (isinstance(e["start"], int) and isinstance(e["end"], int)):
                raise ValueError(f"Bad entity span in sample #{i}: start/end must be int")
            if not (0 <= e["start"] < e["end"] <= len(item["text"])):
                raise ValueError(
                    f"Bad entity span in sample #{i}: [{e['start']}, {e['end']}) out of range for text length {len(item['text'])}"
                )

    return data


def generate_synthetic_samples(n: int = 200, seed: int = 42) -> List[Dict[str, Any]]:
    rnd = random.Random(seed)

    projects = ["ЖК Сосны", "ЖК Олимп", "ЖК Панорама", "проект Северный", "ЖК Речной", "ЖК Сити-Парк"]
    liters = ["литер А", "литер Б", "корпус 1", "корпус 2", "корпус А", "секция 3"]
    rooms_phr = ["студия", "1-комнатная", "2-комнатная", "3-комнатная", "4-комнатная", "трехкомнатная", "двухкомнатная"]
    floor_phr = ["этаж {a}-{b}", "этаж от {a}", "не выше {b} этажа", "с {a} по {b} этаж"]
    sq_phr = ["{a}-{b} м2", "от {a} кв м", "до {b} м²", "{a} м2"]
    cost_phr = ["до {b} млн", "от {a} млн", "{a} {b}00 000 руб", "{a}.{b} млн", "бюджет {a} млн"]

    samples = []
    for _ in range(n):
        r = rnd.choice(rooms_phr)
        fa = rnd.randint(1, 15)
        fb = rnd.randint(fa, 25)
        sqa = rnd.randint(20, 90)
        sqb = rnd.randint(sqa, sqa + rnd.randint(5, 40))
        ca = rnd.randint(2, 20)
        cb = rnd.randint(ca, ca + rnd.randint(1, 10))

        project = rnd.choice(projects)
        liter = rnd.choice(liters)
        floor = rnd.choice(floor_phr).format(a=fa, b=fb)
        sq = rnd.choice(sq_phr).format(a=sqa, b=sqb)
        cost = rnd.choice(cost_phr).format(a=ca, b=cb)

        text = f"ищу {r} {sq}, {cost}, {project}, {liter}, {floor}"

        entities = []

        def add_span(sub: str, label: str):
            start = text.lower().index(sub.lower())
            end = start + len(sub)
            entities.append({"start": start, "end": end, "label": label})

        add_span(r, "ROOMS")
        add_span(sq, "SQUARE")
        add_span(cost, "COST")
        add_span(project, "PROJECT")
        add_span(liter, "LITER")
        add_span(floor, "FLOOR")

        samples.append({"text": text, "entities": entities})

    return samples


# Tokenization + label alignment
def spans_to_token_labels(
    text: str,
    entities: List[Dict[str, Any]],
    tokenizer,
    max_length: int = 128,
) -> Dict[str, Any]:
    # Convert char spans -> BIO token labels using offset mapping.
    enc = tokenizer(
        text,
        truncation=True,
        max_length=max_length,
        return_offsets_mapping=True,
    )
    offsets = enc["offset_mapping"]
    labels = ["O"] * len(offsets)

    # Sort entities by start
    ents = sorted(entities, key=lambda x: (x["start"], x["end"]))

    for e in ents:
        etype = e["label"]
        if etype not in ENTITY_TYPES:
            raise ValueError(f"Unknown entity label: {etype}. Must be in {ENTITY_TYPES}")
        start, end = e["start"], e["end"]
        began = False
        for i, (os_, oe_) in enumerate(offsets):
            if os_ == oe_:
                continue
            # token overlap with entity span?
            overlap = not (oe_ <= start or os_ >= end)
            if overlap:
                if not began:
                    labels[i] = f"B-{etype}"
                    began = True
                else:
                    labels[i] = f"I-{etype}"

    enc.pop("offset_mapping")
    enc["labels"] = [LABEL2ID[l] for l in labels]
    return enc


# Metrics
def compute_metrics_fn(eval_pred):
    import evaluate as hf_evaluate
    seqeval = hf_evaluate.load("seqeval")

    logits, labels = eval_pred
    preds = np.argmax(logits, axis=-1)

    true_labels, true_preds = [], []
    for p, y in zip(preds, labels):
        tl, tp = [], []
        for pi, yi in zip(p, y):
            if yi == -100:
                continue
            tl.append(ID2LABEL[int(yi)])
            tp.append(ID2LABEL[int(pi)])
        true_labels.append(tl)
        true_preds.append(tp)

    res = seqeval.compute(predictions=true_preds, references=true_labels)
    return {
        "precision": res.get("overall_precision", 0.0),
        "recall": res.get("overall_recall", 0.0),
        "f1": res.get("overall_f1", 0.0),
        "accuracy": res.get("overall_accuracy", 0.0),
    }


# Inference: NER decoding
@dataclass
class PredEntity:
    label: str
    start: int
    end: int
    text: str
    score: float


def ner_predict_entities(text: str, tokenizer, model, device: str = "cpu", max_length: int = 192) -> List[PredEntity]:
    import torch

    model.eval()
    model.to(device)

    enc = tokenizer(
        text,
        return_offsets_mapping=True,
        return_tensors="pt",
        truncation=True,
        max_length=max_length,
    )
    offsets = enc.pop("offset_mapping")[0].tolist()
    enc = {k: v.to(device) for k, v in enc.items()}

    with torch.no_grad():
        out = model(**enc)
        logits = out.logits[0]
        probs = torch.softmax(logits, dim=-1)
        pred_ids = torch.argmax(probs, dim=-1).cpu().tolist()
        pred_scores = torch.max(probs, dim=-1).values.cpu().tolist()

    # Build spans from BIO
    entities: List[PredEntity] = []
    cur_label = None
    cur_start = None
    cur_end = None
    cur_scores = []

    for i, (pid, (os_, oe_), sc) in enumerate(zip(pred_ids, offsets, pred_scores)):
        if os_ == oe_:
            continue
        lab = ID2LABEL[pid]
        if lab == "O":
            if cur_label is not None:
                entities.append(_finalize_entity(text, cur_label, cur_start, cur_end, cur_scores))
                cur_label = None
                cur_start = cur_end = None
                cur_scores = []
            continue

        bio, etype = lab.split("-", 1)
        if bio == "B" or (cur_label != etype):
            if cur_label is not None:
                entities.append(_finalize_entity(text, cur_label, cur_start, cur_end, cur_scores))
            cur_label = etype
            cur_start = os_
            cur_end = oe_
            cur_scores = [sc]
        else:
            cur_end = oe_
            cur_scores.append(sc)

    if cur_label is not None:
        entities.append(_finalize_entity(text, cur_label, cur_start, cur_end, cur_scores))

    entities = _merge_close_entities(text, entities, gap=1)
    return entities


def _finalize_entity(text: str, label: str, start: int, end: int, scores: List[float]) -> PredEntity:
    s = float(sum(scores) / max(1, len(scores)))
    return PredEntity(label=label, start=start, end=end, text=text[start:end], score=s)


def _merge_close_entities(text: str, ents: List[PredEntity], gap: int = 1) -> List[PredEntity]:
    if not ents:
        return []
    ents = sorted(ents, key=lambda e: (e.start, e.end))
    merged = [ents[0]]
    for e in ents[1:]:
        last = merged[-1]
        if e.label == last.label and e.start <= last.end + gap:
            new_end = max(last.end, e.end)
            new_text = text[last.start:new_end]
            new_score = (last.score + e.score) / 2.0
            merged[-1] = PredEntity(label=last.label, start=last.start, end=new_end, text=new_text, score=new_score)
        else:
            merged.append(e)
    return merged


UNK = "<UNK>"

ROOMS_WORD2NUM = {
    "студ": 0,
    "одноком": 1,
    "1-ком": 1,
    "1к": 1,
    "двухком": 2,
    "2-ком": 2,
    "2к": 2,
    "трехком": 3,
    "3-ком": 3,
    "3к": 3,
    "четырехком": 4,
    "4-ком": 4,
    "4к": 4,
    "пятиком": 5,
    "5-ком": 5,
    "5к": 5,
}

def normalize_num(s: str) -> Optional[float]:
    s = s.strip().lower()
    s = s.replace(" ", "")
    s = s.replace(",", ".")
    if not re.search(r"\d", s):
        return None
    m = re.findall(r"\d+(?:\.\d+)?", s)
    if not m:
        return None
    try:
        return float(m[0])
    except ValueError:
        return None


def parse_rooms(text: str) -> Tuple[Optional[int], Optional[int]]:
    t = text.lower()
    m = re.search(r"(\d)\s*[- ]?\s*(?:к|комн|комнат)", t)
    if m:
        v = int(m.group(1))
        return v, v

    for k, v in ROOMS_WORD2NUM.items():
        if k in t:
            return v, v

    return None, None


def parse_floor(text: str) -> Tuple[Optional[int], Optional[int]]:
    t = text.lower()

    # 1) этаж 5-12, этажи 5–12, эт 5-12
    m = re.search(r"\b(?:этаж(?:и)?|эт)\.?\s*(\d{1,2})\s*[-–]\s*(\d{1,2})\b", t)
    if m:
        return int(m.group(1)), int(m.group(2))

    # 2) с 5 по 12 этаж
    m = re.search(r"\bс\s*(\d{1,2})\s*по\s*(\d{1,2})\s*(?:этаж|этажа|этажи|эт)\b", t)
    if m:
        return int(m.group(1)), int(m.group(2))

    # 3) этаж от 10" / "эт от 10
    m = re.search(r"\b(?:этаж|эт)\.?\s*от\s*(\d{1,2})\b", t)
    if m:
        return int(m.group(1)), None

    # 4) не выше 8 этажа, до 8 этажа
    m = re.search(r"\b(?:не\s*выше|до)\s*(\d{1,2})\s*этаж(?:а|ей)?\b", t)
    if m:
        return None, int(m.group(1))

    # 5) на 7 этаже
    m = re.search(r"\bна\s*(\d{1,2})\s*этаж(?:е|у)?\b", t)
    if m:
        v = int(m.group(1))
        return v, v

    return None, None


def parse_square(text: str) -> Tuple[Optional[float], Optional[float]]:
    t = text.lower().replace("²", "2")
    # 70-85 м2, 70–85 кв
    m = re.search(r"(\d+(?:[.,]\d+)?)\s*[-–]\s*(\d+(?:[.,]\d+)?)\s*(?:м2|кв\.?\s*м|кв\s*м|кв|м)", t)
    if m:
        a = normalize_num(m.group(1))
        b = normalize_num(m.group(2))
        return a, b

    # от 55 кв м
    m = re.search(r"от\s*(\d+(?:[.,]\d+)?)\s*(?:м2|кв\.?\s*м|кв\s*м|кв|м)", t)
    if m:
        return normalize_num(m.group(1)), None

    # до 30 м2
    m = re.search(r"до\s*(\d+(?:[.,]\d+)?)\s*(?:м2|кв\.?\s*м|кв\s*м|кв|м)", t)
    if m:
        return None, normalize_num(m.group(1))

    # 65 м2
    m = re.search(r"\b(\d+(?:[.,]\d+)?)\s*(?:м2|кв\.?\s*м|кв\s*м)\b", t)
    if m:
        v = normalize_num(m.group(1))
        return v, v

    return None, None


def parse_cost(text: str) -> Tuple[Optional[int], Optional[int]]:
    t = text.lower()

    def to_rub(val: float, mult: int) -> int:
        return int(round(val * mult))

    # до 12 млн, до 4.2млн
    m = re.search(r"до\s*(\d+(?:[.,]\d+)?)\s*(млн|миллион)", t)
    if m:
        b = normalize_num(m.group(1))
        return None, to_rub(b, 1_000_000) if b is not None else (None, None)[1]

    # от 15 млн
    m = re.search(r"от\s*(\d+(?:[.,]\d+)?)\s*(млн|миллион)", t)
    if m:
        a = normalize_num(m.group(1))
        return to_rub(a, 1_000_000) if a is not None else None, None

    # бюджет 9 500 000 руб / 9500000 руб
    m = re.search(r"(\d[\d\s]{5,})\s*(?:руб|р\.?)\b", t)
    if m:
        raw = m.group(1)
        digits = re.sub(r"\s+", "", raw)
        try:
            v = int(digits)
            return v, v
        except ValueError:
            pass

    # 4.2 млн
    m = re.search(r"\b(\d+(?:[.,]\d+)?)\s*млн\b", t)
    if m:
        v = normalize_num(m.group(1))
        return to_rub(v, 1_000_000) if v is not None else None, to_rub(v, 1_000_000) if v is not None else None

    # 4200 тыс
    m = re.search(r"\b(\d+(?:[.,]\d+)?)\s*(тыс|тысяч)\b", t)
    if m:
        v = normalize_num(m.group(1))
        return to_rub(v, 1_000) if v is not None else None, to_rub(v, 1_000) if v is not None else None

    return None, None


def pick_best_text_entity(entities: List[PredEntity], label: str, min_score: float = 0.50) -> Optional[str]:
    cand = [e for e in entities if e.label == label and e.score >= min_score]
    if not cand:
        return None
    cand.sort(key=lambda e: ((e.end - e.start), e.score), reverse=True)
    return cand[0].text.strip(" ,\"'“”«»")


def build_schema_from_text(text: str, entities: List[PredEntity]) -> Dict[str, Any]:
    project = pick_best_text_entity(entities, "PROJECT")
    liter = pick_best_text_entity(entities, "LITER")

    floor_min, floor_max = parse_floor(text)
    rooms_min, rooms_max = parse_rooms(text)
    sq_min, sq_max = parse_square(text)
    cost_min, cost_max = parse_cost(text)

    for e in entities:
        if e.label == "FLOOR" and (floor_min is None and floor_max is None):
            floor_min, floor_max = parse_floor(e.text)
        elif e.label == "ROOMS" and (rooms_min is None and rooms_max is None):
            rooms_min, rooms_max = parse_rooms(e.text)
        elif e.label == "SQUARE" and (sq_min is None and sq_max is None):
            sq_min, sq_max = parse_square(e.text)
        elif e.label == "COST" and (cost_min is None and cost_max is None):
            cost_min, cost_max = parse_cost(e.text)

    def to_out(v):
        return UNK if v is None else v

    return {
        "project_name": project if project else UNK,
        "building_liter": liter if liter else UNK,
        "floor_min": to_out(floor_min),
        "floor_max": to_out(floor_max),
        "rooms_amount_min": to_out(rooms_min),
        "rooms_amount_max": to_out(rooms_max),
        "square_min": to_out(sq_min),
        "square_max": to_out(sq_max),
        "cost_min": to_out(cost_min),
        "cost_max": to_out(cost_max),
    }


# Train
import inspect

def make_training_args(**kwargs):
    """
    Создаёт TrainingArguments, подбирая имена параметров под текущую версию transformers.
    Например, evaluation_strategy vs eval_strategy.
    """
    sig = inspect.signature(TrainingArguments.__init__)
    accepted = set(sig.parameters.keys())

    aliases = {
        "evaluation_strategy": ["eval_strategy"],
        "save_strategy": ["save_strategy"],
        "logging_strategy": ["logging_strategy"],
        "report_to": ["report_to"],
    }

    for main_name, alt_names in aliases.items():
        if main_name in kwargs and main_name not in accepted:
            for alt in alt_names:
                if alt in accepted:
                    kwargs[alt] = kwargs.pop(main_name)
                    break

    filtered = {k: v for k, v in kwargs.items() if k in accepted}
    return TrainingArguments(**filtered)


def train(
    out_dir: str,
    base_model: str = "cointegrated/rubert-tiny",
    seed: int = 42,
    lr: float = 3e-5,
    epochs: int = 6,
    batch_size: int = 16,
    max_length: int = 128,
    train_json: str = DEFAULT_TRAIN_JSON,   # по умолчанию train_data.json
):
    set_seed(seed)

    # 1) Загружаем training samples из файла
    if not train_json:
        raise ValueError("train_json is empty. Provide path to training JSON.")
    if not os.path.exists(train_json):
        raise FileNotFoundError(
            f"Training JSON file not found: '{train_json}'. "
            f"Create it (default name: {DEFAULT_TRAIN_JSON}) or pass --train_json PATH"
        )

    samples = load_training_samples(train_json)
    if len(samples) < 2:
        raise ValueError(f"Training JSON has too few samples: {len(samples)}. Add more examples.")

    # 2) Split train/eval
    random.Random(seed).shuffle(samples)
    split = int(len(samples) * 0.9)
    train_s = samples[:split]
    eval_s = samples[split:] if split < len(samples) else samples[-max(1, len(samples)//10):]

    tokenizer = AutoTokenizer.from_pretrained(base_model)

    def map_fn(ex):
        return spans_to_token_labels(ex["text"], ex["entities"], tokenizer, max_length=max_length)

    ds_train = Dataset.from_list(train_s).map(map_fn, remove_columns=["text", "entities"])
    ds_eval = Dataset.from_list(eval_s).map(map_fn, remove_columns=["text", "entities"])

    model = AutoModelForTokenClassification.from_pretrained(
        base_model,
        num_labels=len(LABELS),
        id2label=ID2LABEL,
        label2id=LABEL2ID,
    )

    args = make_training_args(
        output_dir=out_dir,
        learning_rate=lr,
        num_train_epochs=epochs,
        per_device_train_batch_size=batch_size,
        per_device_eval_batch_size=batch_size,

        evaluation_strategy="epoch",
        save_strategy="epoch",
        logging_strategy="steps",
        logging_steps=20,

        save_total_limit=2,
        load_best_model_at_end=True,
        metric_for_best_model="f1",
        greater_is_better=True,

        fp16=False,  # CUDA
        report_to=[],
    )

    collator = DataCollatorForTokenClassification(tokenizer)

    trainer = Trainer(
        model=model,
        args=args,
        train_dataset=ds_train,
        eval_dataset=ds_eval,
        data_collator=collator,
        tokenizer=tokenizer,
        compute_metrics=compute_metrics_fn,
    )

    trainer.train()
    trainer.save_model(out_dir)
    tokenizer.save_pretrained(out_dir)

    with open(os.path.join(out_dir, "labels.json"), "w", encoding="utf-8") as f:
        json.dump({"labels": LABELS, "entity_types": ENTITY_TYPES}, f, ensure_ascii=False, indent=2)

    print(f"[OK] Saved model to: {out_dir}")
    print(f"[OK] Trained on: {train_json} (train={len(train_s)}, eval={len(eval_s)})")


# Predict (PyTorch)
def predict(model_dir: str, text: str, device: str = "cpu"):
    tokenizer = AutoTokenizer.from_pretrained(model_dir)
    model = AutoModelForTokenClassification.from_pretrained(model_dir)

    entities = ner_predict_entities(text, tokenizer, model, device=device)
    schema = build_schema_from_text(text, entities)

    print(json.dumps(schema, ensure_ascii=False, indent=2))


# ONNX export + ORT inference (optional)
def export_onnx(model_dir: str, onnx_dir: str):
    from optimum.exporters.onnx import main_export

    os.makedirs(onnx_dir, exist_ok=True)

    main_export(
        model_name_or_path=model_dir,
        output=onnx_dir,
        task="token-classification",
    )

    print(f"[OK] Exported ONNX to: {onnx_dir}")


def predict_onnx(onnx_dir: str, text: str, device: str = "cpu"):
    from optimum.onnxruntime import ORTModelForTokenClassification
    import torch

    tokenizer = AutoTokenizer.from_pretrained(onnx_dir)
    model = ORTModelForTokenClassification.from_pretrained(onnx_dir)

    enc = tokenizer(
        text,
        return_offsets_mapping=True,
        return_tensors="pt",
        truncation=True,
        max_length=192,
    )
    offsets = enc.pop("offset_mapping")[0].tolist()

    with torch.no_grad():
        out = model(**enc)
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
                entities.append(_finalize_entity(text, cur_label, cur_start, cur_end, cur_scores))
                cur_label = None
                cur_start = cur_end = None
                cur_scores = []
            continue

        bio, etype = lab.split("-", 1)
        if bio == "B" or (cur_label != etype):
            if cur_label is not None:
                entities.append(_finalize_entity(text, cur_label, cur_start, cur_end, cur_scores))
            cur_label = etype
            cur_start = os_
            cur_end = oe_
            cur_scores = [sc]
        else:
            cur_end = oe_
            cur_scores.append(sc)

    if cur_label is not None:
        entities.append(_finalize_entity(text, cur_label, cur_start, cur_end, cur_scores))

    entities = _merge_close_entities(text, entities, gap=1)
    schema = build_schema_from_text(text, entities)

    print(json.dumps(schema, ensure_ascii=False, indent=2))


# CLI
def main():
    ap = argparse.ArgumentParser()
    sub = ap.add_subparsers(dest="cmd", required=True)

    ap_train = sub.add_parser("train")
    ap_train.add_argument("--out_dir", required=True)
    ap_train.add_argument("--base_model", default="cointegrated/rubert-tiny")
    ap_train.add_argument("--seed", type=int, default=42)
    ap_train.add_argument("--lr", type=float, default=3e-5)
    ap_train.add_argument("--epochs", type=int, default=6)
    ap_train.add_argument("--batch_size", type=int, default=16)
    ap_train.add_argument("--max_length", type=int, default=128)

    ap_pred = sub.add_parser("predict")
    ap_pred.add_argument("--model_dir", required=True)
    ap_pred.add_argument("--text", required=True)
    ap_pred.add_argument("--device", default="cpu")

    ap_export = sub.add_parser("export-onnx")
    ap_export.add_argument("--model_dir", required=True)
    ap_export.add_argument("--onnx_dir", required=True)

    ap_pred_onnx = sub.add_parser("predict-onnx")
    ap_pred_onnx.add_argument("--onnx_dir", required=True)
    ap_pred_onnx.add_argument("--text", required=True)
    ap_pred_onnx.add_argument("--device", default="cpu")

    ap_train.add_argument("--train_json", default=DEFAULT_TRAIN_JSON, help="Path to training JSON")

    args = ap.parse_args()

    if args.cmd == "train":
        train(
            out_dir=args.out_dir,
            base_model=args.base_model,
            seed=args.seed,
            lr=args.lr,
            epochs=args.epochs,
            batch_size=args.batch_size,
            max_length=args.max_length,
            train_json=args.train_json,
        )
    elif args.cmd == "predict":
        predict(args.model_dir, args.text, device=args.device)
    elif args.cmd == "export-onnx":
        export_onnx(args.model_dir, args.onnx_dir)
    elif args.cmd == "predict-onnx":
        predict_onnx(args.onnx_dir, args.text, device=args.device)
    else:
        raise RuntimeError("Unknown command")


if __name__ == "__main__":
    main()

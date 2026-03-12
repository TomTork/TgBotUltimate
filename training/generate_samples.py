import json
import random
from typing import Dict, List, Any, Tuple

ENTITY_TYPES = ["PROJECT", "LITER", "FLOOR", "ROOMS", "SQUARE", "COST"]

def ci_find_span(text: str, substr: str) -> Tuple[int, int]:
    lt = text.lower()
    ls = substr.lower()
    i = lt.find(ls)
    if i < 0:
        raise ValueError(f"Substring not found: {substr!r} in text: {text!r}")
    return i, i + len(substr)

def try_add_entity(entities: List[Dict[str, Any]], text: str, substr: str, label: str) -> bool:
    if not substr:
        return False
    lt = text.lower()
    ls = substr.lower()
    i = lt.find(ls)
    if i < 0:
        return False
    entities.append({"start": i, "end": i + len(substr), "label": label})
    return True

def maybe_typo(rnd: random.Random, s: str, p: float = 0.30) -> str:
    if rnd.random() > p:
        return s
    ops = [
        lambda x: x.replace("ищу", "ишу"),
        lambda x: x.replace("квартиру", "квартриру"),
        lambda x: x.replace("нужна", "нужн"),
        lambda x: x.replace(", ", ","),
        lambda x: x.replace("  ", " "),
        lambda x: x.replace(" м2", "м2"),
        lambda x: x.replace(" до ", " до"),
        lambda x: x.replace(" от ", " от"),
        lambda x: x.replace("этаж ", "этаж"),
        lambda x: x.replace("жк ", "жк"),
    ]
    return rnd.choice(ops)(s)

def rooms_phrase(rnd: random.Random) -> str:
    options = [
        "студия", "студ", "студию",

        "1к", "1 к", "1-комн", "1-комнатная", "1-комнатную",
        "однушка", "однокомнатная", "однокомнатную", "одноком",

        "2к", "2 к", "2-комн", "2-комнатная", "2-комнатную",
        "двушка", "двухкомнатная", "двухкомнатную", "двухком",

        "3к", "3 к", "3-комн", "3-комнатная", "3-комнатную",
        "трёшка", "трехкомнатная", "трехкомнатную", "трехком",

        "4к", "4 к", "4-комн", "4-комнатная", "4-комнатную",
        "четырехкомнатная", "четырехкомнатную", "четырехком",
    ]
    return rnd.choice(options)

def square_phrase(rnd: random.Random) -> str:
    a = rnd.randint(18, 95)
    b = rnd.randint(a, a + rnd.randint(3, 45))
    unit = rnd.choice(["м2", "м²", "кв м", "кв.м", "квм"])
    dash = rnd.choice(["-", "–", "—"])
    mode = rnd.choice(["range", "min", "max", "exact", "broken_range"])
    if mode == "range":
        return f"{a}{dash}{b} {unit}"
    if mode == "broken_range":
        return f"{a} {b}{unit}"         # "70 85м2"
    if mode == "min":
        return f"от {a} {unit}"
    if mode == "max":
        return f"до {b}{unit}" if rnd.random() < 0.5 else f"до {b} {unit}"
    return f"{a} {unit}"

def floor_phrase(rnd: random.Random) -> str:
    a = rnd.randint(1, 20)
    b = rnd.randint(a, 30)
    dash = rnd.choice(["-", "–", "—"])
    kind = rnd.choice(["этаж", "этажи", "эт", "эт."])
    mode = rnd.choice(["range", "min", "max", "exact", "no_space"])
    if mode == "range":
        return f"{kind} {a}{dash}{b}"
    if mode == "no_space":
        return f"{kind}{a}{dash}{b}"
    if mode == "min":
        return f"{kind} от {a}"
    if mode == "max":
        return f"до {b} этажа" if rnd.random() < 0.5 else f"не выше {b} этажа"
    return f"на {a} этаже"

def cost_phrase(rnd: random.Random) -> str:
    mln = rnd.uniform(2.0, 25.0)
    mln = round(mln, 1) if rnd.random() < 0.5 else round(mln)
    rub = rnd.randint(2_500_000, 25_000_000)
    mode = rnd.choice(["max_mln", "min_mln", "exact_mln", "exact_rub", "budget_mln", "compact"])
    if mode == "max_mln":
        return f"до {mln} млн"
    if mode == "min_mln":
        return f"от {mln} млн"
    if mode == "budget_mln":
        return f"бюджет {mln} млн"
    if mode == "compact":
        prefix = rnd.choice(["до", "от"])
        val = str(mln).replace(".0", "")
        return f"{prefix}{val}млн"
    if mode == "exact_mln":
        return f"{mln} млн"
    return f"{rub:,}".replace(",", " ") + rnd.choice([" руб", "р", " руб."])

def project_phrase(rnd: random.Random) -> str:
    names = [
        "ЖК Сосны", "ЖК Олимп", "ЖК Панорама", "ЖК Речной", "ЖК Сити Парк",
        "жк северный", "жк южный", "жк лесной", "жк маяк", "жк аврора",
        "проект Северный", "проект Ривьера", "комплекс Сады",
        "ЖК 'Панорама'", "ЖК «Горизонт»", "ЖК \"Берег\""
    ]
    s = rnd.choice(names)
    if rnd.random() < 0.25:
        s = s.replace("ЖК", "жк")
    if rnd.random() < 0.15:
        s = s.replace(" ", "")
    return s

def liter_phrase(rnd: random.Random) -> str:
    lit = rnd.choice(["А", "Б", "В", "Г", "1", "2", "3"])
    kind = rnd.choice(["литер", "лит", "л", "корпус", "корп", "секция"])
    sep = rnd.choice([" ", "", " №", " #"])
    return f"{kind}{sep}{lit}"

def make_text(rnd: random.Random) -> Tuple[str, List[Dict[str, Any]]]:
    r = rooms_phrase(rnd)
    sq = square_phrase(rnd)
    fl = floor_phrase(rnd)
    co = cost_phrase(rnd)
    pr = project_phrase(rnd)
    li = liter_phrase(rnd)

    style = rnd.choice(["full", "semi", "short", "ultra_short", "scrambled"])
    parts: List[str] = []

    include_sq = True
    include_fl = True
    include_li = True

    if style == "full":
        parts = [
            rnd.choice(["ищу", "нужна", "надо", "хочу"]),
            rnd.choice(["квартира", "квартиру", "кв"]),
            r,
            sq + ",",
            co + ",",
            pr + ",",
            li + ",",
            fl,
        ]
    elif style == "semi":
        parts = [rnd.choice(["ищу", "хочу", "нужна"]), r, pr + ",", co + ",", sq + ",", fl, li]
    elif style == "short":
        include_li = rnd.random() < 0.35
        include_sq = True
        include_fl = rnd.random() < 0.75
        parts = [r, pr, co]
        if include_sq:
            parts.append(sq)
        if include_fl:
            parts.append(fl)
        if include_li:
            parts.append(li)
        text = " ".join(parts)
        if rnd.random() < 0.6:
            text = text.replace(" ", ", ", 1)
        parts = [text]
    elif style == "ultra_short":
        include_sq = rnd.random() < 0.45
        include_fl = rnd.random() < 0.35
        include_li = rnd.random() < 0.25
        parts = [r, co, pr]
        if include_sq:
            parts.append(sq)
        if include_fl:
            parts.append(fl)
        if include_li:
            parts.append(li)
    else:
        include_sq = rnd.random() < 0.85
        include_fl = rnd.random() < 0.85
        include_li = rnd.random() < 0.75
        parts = [pr, r, co]
        if include_sq:
            parts.append(sq)
        if include_fl:
            parts.append(fl)
        if include_li:
            parts.append(li)
        rnd.shuffle(parts)

    text = " ".join(parts)
    text = maybe_typo(rnd, text, p=0.35)

    entities: List[Dict[str, Any]] = []
    try_add_entity(entities, text, pr, "PROJECT")
    try_add_entity(entities, text, r, "ROOMS")
    try_add_entity(entities, text, co, "COST")
    if include_sq:
        try_add_entity(entities, text, sq, "SQUARE")
    if include_fl:
        try_add_entity(entities, text, fl, "FLOOR")
    if include_li:
        try_add_entity(entities, text, li, "LITER")

    if len(entities) < 3:
        return make_text(rnd)

    return text, entities

def main():
    import argparse
    ap = argparse.ArgumentParser()
    ap.add_argument("--out", default="train_data.json")
    ap.add_argument("--n", type=int, default=2000)
    ap.add_argument("--seed", type=int, default=42)
    args = ap.parse_args()

    rnd = random.Random(args.seed)

    samples = []
    for _ in range(args.n):
        text, entities = make_text(rnd)

        for e in entities:
            if e["label"] not in ENTITY_TYPES:
                raise ValueError("Bad label")
            if not (0 <= e["start"] < e["end"] <= len(text)):
                raise ValueError("Bad span range")
            if not text[e["start"]:e["end"]]:
                raise ValueError("Empty span")

        samples.append({"text": text, "entities": entities})

    with open(args.out, "w", encoding="utf-8") as f:
        json.dump(samples, f, ensure_ascii=False, indent=2)

    print(f"[OK] Wrote {len(samples)} samples to {args.out}")

if __name__ == "__main__":
    main()

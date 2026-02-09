from dataclasses import dataclass
from typing import Dict, Any


@dataclass
class LeadScoreResult:
    score: float
    model_version: str


def simple_lead_score(features: Dict[str, Any]) -> LeadScoreResult:
    """
    Placeholder scoring. Replace with a real model later.
    """
    # Example heuristic: bias toward leads with email + phone + recent activity.
    score = 0.2
    if features.get("has_email"):
        score += 0.2
    if features.get("has_phone"):
        score += 0.2
    if features.get("notes_count", 0) > 0:
        score += 0.2
    if features.get("recent_activity"):
        score += 0.2

    return LeadScoreResult(score=min(score, 0.95), model_version="v0.1-heuristic")


@dataclass
class ChurnResult:
    risk_score: float
    model_version: str


def simple_churn_score(features: Dict[str, Any]) -> ChurnResult:
    """
    Placeholder churn scoring.
    Higher score means higher churn risk.
    """
    risk = 0.1
    if not features.get("recent_activity"):
        risk += 0.4
    if features.get("notes_count", 0) == 0:
        risk += 0.2
    if features.get("closed_won", 0) == 0 and features.get("total_deals", 0) > 0:
        risk += 0.2
    return ChurnResult(risk_score=min(risk, 0.95), model_version="v0.1-heuristic")


@dataclass
class CLVResult:
    clv: float
    model_version: str


def simple_clv_score(features: Dict[str, Any]) -> CLVResult:
    """
    Placeholder CLV scoring: uses total deal value with a discount factor.
    """
    total_value = float(features.get("total_deal_value", 0.0))
    clv = total_value * 0.85
    return CLVResult(clv=clv, model_version="v0.1-heuristic")

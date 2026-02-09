from fastapi import FastAPI, HTTPException

from .config import load_config
from .warehouse import get_client
from .models import (
    LeadScoreRequest,
    LeadScoreResponse,
    ChurnRequest,
    ChurnResponse,
    CLVRequest,
    CLVResponse,
)
from .lead_scoring import simple_lead_score, simple_churn_score, simple_clv_score

app = FastAPI(title="AI Service", version="0.1.0")

cfg = load_config()
print("ye dekho")
print(cfg)
client = get_client(cfg.warehouse)


@app.get("/health")
def health():
    return {"status": "ok"}


@app.post("/v1/lead-score", response_model=LeadScoreResponse)
def lead_score(req: LeadScoreRequest):
    # Minimal feature extraction from warehouse for lead/contact.
    query = """
        SELECT
            id,
            email != '' AS has_email,
            phone != '' AS has_phone,
            notes_count,
            last_activity_at >= now() - INTERVAL 30 DAY AS recent_activity
        FROM contacts_features
        WHERE id = {lead_id:String}
        LIMIT 1
    """

    rows = client.query(query, parameters={"lead_id": req.lead_id}).result_rows
    if not rows:
        raise HTTPException(status_code=404, detail="lead not found in warehouse")

    _, has_email, has_phone, notes_count, recent_activity = rows[0]
    features = {
        "has_email": bool(has_email),
        "has_phone": bool(has_phone),
        "notes_count": int(notes_count),
        "recent_activity": bool(recent_activity),
    }
    result = simple_lead_score(features)
    return LeadScoreResponse(lead_id=req.lead_id, score=result.score, model_version=result.model_version)


@app.post("/v1/churn-score", response_model=ChurnResponse)
def churn_score(req: ChurnRequest):
    query = """
        SELECT
            contact_id,
            last_activity_at >= now() - INTERVAL 30 DAY AS recent_activity,
            notes_count,
            total_deals,
            closed_won
        FROM churn_features
        WHERE contact_id = {contact_id:String}
        LIMIT 1
    """

    rows = client.query(query, parameters={"contact_id": req.contact_id}).result_rows
    if not rows:
        raise HTTPException(status_code=404, detail="contact not found in warehouse")

    _, recent_activity, notes_count, total_deals, closed_won = rows[0]
    features = {
        "recent_activity": bool(recent_activity),
        "notes_count": int(notes_count),
        "total_deals": int(total_deals),
        "closed_won": int(closed_won),
    }
    result = simple_churn_score(features)
    return ChurnResponse(contact_id=req.contact_id, risk_score=result.risk_score, model_version=result.model_version)


@app.post("/v1/clv", response_model=CLVResponse)
def clv(req: CLVRequest):
    query = """
        SELECT
            contact_id,
            total_deal_value,
            avg_deal_value,
            total_deals,
            closed_won
        FROM clv_features
        WHERE contact_id = {contact_id:String}
        LIMIT 1
    """

    rows = client.query(query, parameters={"contact_id": req.contact_id}).result_rows
    if not rows:
        raise HTTPException(status_code=404, detail="contact not found in warehouse")

    _, total_deal_value, avg_deal_value, total_deals, closed_won = rows[0]
    features = {
        "total_deal_value": float(total_deal_value),
        "avg_deal_value": float(avg_deal_value),
        "total_deals": int(total_deals),
        "closed_won": int(closed_won),
    }
    result = simple_clv_score(features)
    return CLVResponse(contact_id=req.contact_id, clv=result.clv, model_version=result.model_version)

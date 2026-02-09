from pydantic import BaseModel


class LeadScoreRequest(BaseModel):
    lead_id: str


class LeadScoreResponse(BaseModel):
    lead_id: str
    score: float
    model_version: str


class ChurnRequest(BaseModel):
    contact_id: str


class ChurnResponse(BaseModel):
    contact_id: str
    risk_score: float
    model_version: str


class CLVRequest(BaseModel):
    contact_id: str


class CLVResponse(BaseModel):
    contact_id: str
    clv: float
    model_version: str

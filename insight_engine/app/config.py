from pydantic import BaseModel
import tomllib


class AppConfig(BaseModel):
    env: str = "dev"
    host: str = "0.0.0.0"
    port: int = 8090


class WarehouseConfig(BaseModel):
    host: str
    port: int
    database: str
    user: str
    password: str


class Config(BaseModel):
    app: AppConfig
    warehouse: WarehouseConfig


def load_config(path: str = "./config/default.toml") -> Config:
    with open(path, "rb") as f:
        data = tomllib.load(f)
    app_raw = data.get("App", {})
    wh_raw = data.get("Warehouse", {})
    return Config(
        app=AppConfig(
            env=app_raw.get("Env", "dev"),
            host=app_raw.get("Host", "0.0.0.0"),
            port=int(app_raw.get("Port", 8090)),
        ),
        warehouse=WarehouseConfig(
            host=wh_raw.get("Host", "127.0.0.1"),
            port=int(wh_raw.get("Port", 8123)),
            database=wh_raw.get("Database", "crm_warehouse"),
            user=wh_raw.get("User", "default"),
            password=wh_raw.get("Password", ""),
        ),
    )

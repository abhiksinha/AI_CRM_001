import clickhouse_connect
from .config import WarehouseConfig


def get_client(cfg: WarehouseConfig):
    return clickhouse_connect.get_client(
        host=cfg.host,
        port=cfg.port,
        username=cfg.user,
        password=cfg.password,
        database=cfg.database,
    )

from sqlalchemy import create_engine
from sqlalchemy.orm import sessionmaker
from config.secret_config import config
from db.models import Base

DB_URI = config.get("DatabaseURI", "sqlite:///./data/events.db")

engine = create_engine(DB_URI, echo=False)
SessionLocal = sessionmaker(bind=engine)

def init_db():
    Base.metadata.create_all(bind=engine)
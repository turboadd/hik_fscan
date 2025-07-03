from sqlalchemy import Column, Integer, String, DateTime
from sqlalchemy.ext.declarative import declarative_base
from datetime import datetime

Base = declarative_base()

class Event(Base):
    __tablename__ = "events"
    id = Column(Integer, primary_key=True)
    site_id = Column(Integer)
    ip = Column(String)
    device_id = Column(Integer)
    client_id = Column(Integer)
    event_time = Column(DateTime)
    received_at = Column(DateTime, default=datetime.utcnow)
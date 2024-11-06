from fastapi import FastAPI, HTTPException
from fastapi.middleware.cors import CORSMiddleware
from pydantic import BaseModel
import os
import mysql.connector
from dotenv import load_dotenv

# Cargar variables de entorno desde un archivo .env
load_dotenv()

# Configurar FastAPI
app = FastAPI()

# Configura CORS
origins = [
    "http://localhost:3000",  # React frontend running locally
    "http://frontend:3000",
    "http://127.0.0.1:3000"
]
app.add_middleware(
    CORSMiddleware,
    allow_origins=origins,
    allow_credentials=True,
    allow_methods=["GET", "POST", "PUT", "DELETE", "OPTIONS"],
    allow_headers=["*"],
    expose_headers=["*"],
    max_age=3600,
)

# Estructura para el modelo de Imagen
class Image(BaseModel):
    id: int
    link: str

# Configuración de la conexión a MySQL
def get_db_connection():
    return mysql.connector.connect(
        host=os.getenv("MYSQL_HOST", "mysql"),
        user=os.getenv("MYSQL_USER", "root"),
        password=os.getenv("MYSQL_PASSWORD", ""),
        database=os.getenv("MYSQL_DATABASE", "test")
    )

@app.get("/")
def read_root():
    return {"message": "API is running"}

# Endpoint para obtener las imágenes
@app.get("/images", response_model=list[Image])
def get_images():
    try:
        conn = get_db_connection()
        cursor = conn.cursor(dictionary=True)
        cursor.execute("SELECT id, link FROM images")
        images = cursor.fetchall()
        cursor.close()
        conn.close()
        return images
    except mysql.connector.Error as err:
        raise HTTPException(status_code=500, detail=f"Database error: {err}")

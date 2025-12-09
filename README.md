# 💡 UNDP RAG System for Afghanistan Value Chain Integration Tool

This project implements a **Retrieval-Augmented Generation (RAG)** system to support UNDP teams working on MSMEs in Least Developed Countries (LDCs), with a focus on Afghanistan. It enables document-based question answering, summarization, decision support, and metadata extraction using Azure OpenAI and ChromaDB.

---

## 🚀 Project Goals

- Provide UNDP staff with AI-powered access to internal reports and policy briefs.
- Enable fast and accurate retrieval of insights from large document collections.
- Support programme design, research, and capacity building efforts.

---

## 🧠 Technologies Used

- **Python (Django REST Framework)** – Backend API service  
- **Next.js** – Frontend (if applicable)  
- **Azure OpenAI** – GPT models + embeddings  
- **ChromaDB** – Vector store for RAG retrieval  
- **PostgreSQL** – Metadata + app storage  
- **Docker & Docker Compose** – Full environment orchestration  
- **Nginx** – Reverse proxy for the backend  
- **LangChain** – Framework for building LLM-powered applications

---

## 🧩 Prompt Categories

The system supports three types of queries:

1. **Research / Analyses**  
   _e.g., "What are the key impacts of the PVPV Morality Law on women-led MSMEs in Afghanistan?"_

2. **Programme Design**  
   _e.g., "How is UNDP supporting women-led MSMEs in Afghanistan through programme design?"_

3. **Capacity Building**  
   _e.g., "What capacity building initiatives has UNDP implemented to support Afghan women entrepreneurs?"_

---

# ⚙️ Running the Project (Backend + Frontend)

Below is the updated **step-by-step guide** based on all issues previously encountered during setup.

---

## ✅ **1. Prerequisites**

Before running the project, ensure you have:

### **System**
- Docker Desktop (required)
- Docker Compose v2+
- Git

### **Azure**
You must have:
- `AZURE_OPENAI_ENDPOINT`
- `AZURE_OPENAI_API_KEY`
- `AZURE_OPENAI_API_VERSION`
- `AZURE_OPENAI_DEPLOYMENT_NAME` (Chat model)
- `AZURE_EMBEDDINGS_DEPLOYMENT` (Embedding model)

These should be placed in your `.env` file.

---

## 🔧 **2. Environment Setup**

Copy `.env.example` → `.env` and fill with your Azure credentials:

```bash
cp .env.example .env
```

Edit the `.env` file with your Azure credentials:

```ini
AZURE_OPENAI_API_KEY=your-key
AZURE_OPENAI_ENDPOINT=https://your-endpoint.openai.azure.com/
AZURE_OPENAI_API_VERSION=2024-06-01
AZURE_OPENAI_DEPLOYMENT_NAME=gpt-4o-mini
AZURE_EMBEDDINGS_DEPLOYMENT=text-embedding-3-large
```

### **Alternative: Manual Setup (Without Docker)**

If you prefer to run without Docker:

1. Install Python dependencies:
```bash
pip install -r requirements.txt
```

2. Set your Azure credentials:
```bash
export AZURE_OPENAI_API_KEY='your-key'
export AZURE_OPENAI_ENDPOINT='your-endpoint'
export AZURE_OPENAI_DEPLOYMENT_NAME='your-chat-deployment'
export AZURE_EMBEDDINGS_DEPLOYMENT='your-embedding-deployment'
```

3. Add Documents:
Place your PDF files in the `data/` folder.

4. Populate the Vector Database:
```bash
python populate_database.py --reset
```

5. Run a Query:
```bash
python query_data.py "Your question here"
```

---

## 🚀 **3. Start the Entire Stack**

From project root:

```bash
docker compose up --build
```

This will launch:
- Django backend
- Celery worker
- Celery beat
- PostgreSQL
- Nginx
- Frontend (if included)

---

## 🟦 **Backend Usage**

### **Access Django container**
```bash
docker exec -it django bash
```

### **Run migrations manually** (only needed if migrations change)
```bash
python manage.py migrate
```

### **Check URLs inside Django**
(Useful when Nginx routing errors happen)

```bash
python manage.py show_urls
```

---

## ⚠️ **Troubleshooting Guide** (Updated From Your Issues)

These are the most common problems we encountered and the solutions now included in the README.

### ❗ **1. 405 Errors (Method Not Allowed)**

**Cause:**
Using the wrong API route.
`/chat/` does not accept POST directly.

**Fix:**
Use:
```bash
POST /api/chat/
```

Nginx rewrites routes, so testing backend directly must use `/api/chat/`.

### ❗ **2. 401 "Invalid Token" Errors**

**Cause:**
Missing or incorrect `Authorization: Bearer <key>` header.

**Fix:**
Make sure your Postman request includes:
```vbnet
Authorization: Bearer dev-key
```
Or whatever key your backend expects.

### ❗ **3. Django Fails to Start**

`ModuleNotFoundError: No module named 'django_extensions'`

**Cause:**
`django-extensions` was added to `INSTALLED_APPS` but not installed.

**Fix:**
Add it to `requirements.txt` and rebuild:
```bash
docker compose up --build
```

### ❗ **4. Nginx → Django Routing Mismatch**

**Symptoms:**
- `/chat/sessions/...` gives 404
- `/chat/?session_id=1` gives 405

**Fix Checklist:**
1. Confirm backend routes using:
```bash
python manage.py show_urls
```

2. Always test through Nginx at `/api/...`

3. Inside Django container, you can test backend directly via:
```bash
curl http://localhost:8000/chat/
```

### ❗ **5. Environment Variables Not Loaded**

If Django logs show missing Azure credentials:

**Fix:**
Restart all containers:
```bash
docker compose down
docker compose up --build
```

Confirm variables inside Django:
```bash
docker exec -it django bash
env | grep AZURE
```

---

## 📁 **Project Structure (High Level)**

```
backend/
  rag/
    views.py
    services/
      prompt.py
      rag_pipeline.py
    utils/
frontend/
nginx/
docker-compose.yml
.env
```

---

## 🧪 **Testing Guardrails**

To test the guardrail functionality:
```bash
python test_guardrails.py
```

This prints Markdown output you can paste into PRs.

---

## 📄 **License**

This project is for internal UNDP use. Contact the author for reuse or adaptation.

---

## 📞 **Support**

For technical support or questions:
1. Check the troubleshooting guide above
2. Ensure all environment variables are set correctly
3. Verify Azure OpenAI service is accessible from your network
4. For persistent issues, contact the development team with error logs and environment details

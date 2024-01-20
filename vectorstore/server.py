import os
from fastapi import FastAPI
from langchain.llms import Ollama
from langchain.embeddings.ollama import OllamaEmbeddings
from llama_index import load_index_from_storage
from llama_index import ServiceContext, StorageContext, VectorStoreIndex
from llama_index import Document
from pydantic import BaseModel
from typing import List


class Message(BaseModel):
    role: str
    content: str


app = FastAPI()


llm = Ollama(
    base_url="http://localhost:3020",
    model="mistral:instruct",
    temperature=0.01,
    additional_kwargs={"top_p": 1, "max_new_tokens": 256},
)

embed_model = OllamaEmbeddings(model="mistral:instruct")
service_context = ServiceContext.from_defaults(
    llm=llm, embed_model=embed_model
)


@app.post("/api/vectorstore/add")
def add_to_vectorstore(messages: List[Message]):
    markdown_document = "\n\n".join(
        f"## {message.role}\n\n{message.content}"
        for message in messages
    )
    new_doc = Document(text=markdown_document)
    print(new_doc)

    storage_context = StorageContext.from_defaults(persist_dir="./storage")
    if not os.path.exists("./storage"):
        index = VectorStoreIndex.from_documents(
            [new_doc],
            service_context=service_context,
        )
    else:
        index = load_index_from_storage(storage_context)

    index.storage_context.persist()


if __name__ == "__main__":
    import uvicorn
    uvicorn.run(app, port=3030)

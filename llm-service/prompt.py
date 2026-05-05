def prompt_messages(question: str, context: list) -> list:
    system_prompt = (
        "You are an incident investigation copilot.\n"
        "Use only the context below to answer the question.\n"
        "If the context doesn't contain enough information, say so."
    )

    user_prompt = f"Context:\n{context}\n\nQuestion:\n{question}"

    messages = [
        {"role": "system", "content": system_prompt},
        {"role": "user", "content": user_prompt}
    ]

    return messages


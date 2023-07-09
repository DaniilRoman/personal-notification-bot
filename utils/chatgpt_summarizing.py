import openai
import logging


def configure_openai(api_key: str, organization: str):
    openai.api_key = api_key
    openai.organization = organization


def summarize_text(text: str) -> str:
    input_text = f"Give me the main ideas of this tech article: {text[:4000]}"  # model limit
    model = "gpt-3.5-turbo"
    max_tokens = 500

    try:
        # Generate a summary using OpenAI's complete() function
        response = openai.ChatCompletion.create(
            model=model,
            messages=[
                {'role': 'user', 'content': input_text}
            ],
            temperature=0,
            max_tokens=max_tokens
        )

        # Extract the generated summary from the response
        summary = response['choices'][0]['message']['content'].strip()

        return summary
    except Exception as ex:
        logging.warning("Couldn't get summary for the text", ex)
        return ""

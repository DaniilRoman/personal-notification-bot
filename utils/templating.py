from jinja2 import Environment, FileSystemLoader
import os

MAIN_DIR = os.path.dirname(os.path.abspath(__package__))
WWW_DIR = os.path.join(MAIN_DIR, "www")


def render_index_html(data: dict):
    j2_env = Environment(loader=FileSystemLoader(WWW_DIR),
                         trim_blocks=True)
    rendered_html = j2_env.get_template('template.html').render(data)

    index_html = os.path.join(WWW_DIR, "index.html")
    with open(index_html, "w") as index_html:
        index_html.write(rendered_html)

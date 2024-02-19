from flask import Flask, render_template, session, request
from flask_babel import Babel, _
import os
import requests

def create_app():
    app = Flask(__name__)

    # Set a secret key for the session
    app.config['SECRET_KEY'] = os.urandom(24)

    # Configure Flask-Babel
    babel = Babel(app)
    app.config['LANGUAGES'] = ['en', 'es']
    app.config['BABEL_DEFAULT_LOCALE'] = 'es'

    # Replace this URL with your Go GraphQL server endpoint
    GRAPHQL_URL = "http://backend:8000/graphql"

    def fetch_images():
        query = """
        {
            images {
                ID
                Link
            }
        }
        """
        response = requests.post(GRAPHQL_URL, json={'query': query})
        data = response.json()
        return data['data']['images']

    @babel.localeselector
    def get_locale():
        # Use explicit comparison and validate against supported languages
        return session.get('lang') in app.config['LANGUAGES'] and session.get('lang') or 'es'

    @app.before_request
    def before_request():
        if request.endpoint == 'static':
            return  # Skip language setting for static files

        session['lang'] = request.args.get('lang', session.get('lang'))

    @app.route('/')
    def index():
        images = fetch_images()
        return render_template('index.html', images=images, get_locale=get_locale)

    @app.route('/set_language/<lang>', methods=['POST'])
    def set_language(lang):
        if lang in app.config['LANGUAGES']:
            session['lang'] = lang
        return ('', 204)

    return app

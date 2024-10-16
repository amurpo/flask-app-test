from flask import Flask, render_template, session, request
from flask_babel import Babel, _
import os  # Import the 'os' module for generating a secret key
import requests

def get_user_locale():
    return session.get('lang', 'en')

def create_app():
    app = Flask(__name__)

    # Set a secret key for the session from environment variable
    app.config['SECRET_KEY'] = os.getenv('SECRET_KEY', os.urandom(24).hex())  # Use hex to get a string representation

    # Configuración de la cookie de sesión (se configurará más adelante)
    app.config['SESSION_COOKIE_SAMESITE'] = 'Lax'  # Valor predeterminado
    app.config['SESSION_COOKIE_SECURE'] = False  # Valor predeterminado

    # Configure Flask-Babel
    babel = Babel(app)
    app.config['LANGUAGES'] = ['en', 'es', 'tr', 'ko']
    app.config['BABEL_DEFAULT_LOCALE'] = 'es'

    # Ensure locale selection after Flask initialization
    @babel.localeselector
    def get_locale():
        return get_user_locale()

    # Replace this URL with your Go GraphQL server endpoint
    GRAPHQL_URL = "http://backend:8000/graphql"

    def fetch_images():
        query = '''
        {
            images {
                ID
                Link
            }
        }
        '''

        response = requests.post(GRAPHQL_URL, json={'query': query})
        data = response.json()
        return data['data']['images']

    @app.before_request
    def before_request():
        if request.endpoint == 'static':
            return  # Skip language setting for static files

        session['lang'] = request.args.get('lang', session.get('lang', 'en'))

        # Configura la cookie de sesión según el protocolo de la solicitud
        if request.headers.get('X-Forwarded-Proto') == 'https' or request.scheme == 'https':
            app.config['SESSION_COOKIE_SAMESITE'] = 'None'
            app.config['SESSION_COOKIE_SECURE'] = True
        else:
            app.config['SESSION_COOKIE_SAMESITE'] = 'Lax'
            app.config['SESSION_COOKIE_SECURE'] = False

    @app.route('/')
    def index():
        images = fetch_images()
        context = {'get_locale': get_locale, 'gettext': _}  # Include get_locale in the context
        print(_('Hello, I am a bear!'))  # Imprimirá la traducción en la consola de la aplicación
        return render_template('index.html', images=images, **context)

    @app.route('/set_language/<lang>', methods=['POST'])
    def set_language(lang):
        if lang in app.config['LANGUAGES']:
            session['lang'] = lang
        return ('', 204)

    return app


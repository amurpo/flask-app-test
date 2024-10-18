from flask import Flask, render_template, session, request
from flask_babel import Babel, _
import os
import requests

def get_user_locale():
    return session.get('lang', 'en')

def create_app():
    app = Flask(__name__)
    # Set a secret key for the session from environment variable
    app.config['SECRET_KEY'] = os.getenv('SECRET_KEY', os.urandom(24).hex())
    
    # Session cookie configuration
    app.config['SESSION_COOKIE_SAMESITE'] = 'Lax'
    app.config['SESSION_COOKIE_SECURE'] = False
    
    # Configure Flask-Babel
    app.config['LANGUAGES'] = ['en', 'es', 'tr', 'ko']
    app.config['BABEL_DEFAULT_LOCALE'] = 'es'
    babel = Babel(app)
    
    # Set the locale selector function directly
    babel.init_app(app, locale_selector=get_user_locale)

    # GraphQL endpoint
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
        try:
            response = requests.post(GRAPHQL_URL, json={'query': query})
            response.raise_for_status()
            data = response.json()
            return data.get('data', {}).get('images', [])
        except requests.exceptions.RequestException as e:
            print(f"Error fetching images: {e}")
            return []

    @app.before_request
    def before_request():
        if request.endpoint == 'static':
            return
        
        session['lang'] = request.args.get('lang', session.get('lang', 'en'))
        
        is_secure = request.headers.get('X-Forwarded-Proto') == 'https' or request.scheme == 'https'
        app.config['SESSION_COOKIE_SAMESITE'] = 'None' if is_secure else 'Lax'
        app.config['SESSION_COOKIE_SECURE'] = is_secure

    @app.route('/')
    def index():
        images = fetch_images()
        context = {
            'get_locale': get_user_locale,
            'gettext': _
        }
        print(_('Hello, I am a bear!'))
        return render_template('index.html', images=images, **context)

    @app.route('/set_language/<lang>', methods=['POST'])
    def set_language(lang):
        if lang in app.config['LANGUAGES']:
            session['lang'] = lang
        return ('', 204)

    return app

from flask import Flask, render_template

def create_app():
    app = Flask(__name__)

    @app.route('/')
    def hello():
        img_url = 'https://oso-flask.s3.sa-east-1.amazonaws.com/oso.jpg'
        message = "Hola, Soy un oso!"
        return render_template('index.html', message=message, img_url=img_url)

    return app

if __name__ == '__main__':
    app = create_app()
    app.run(debug=True)

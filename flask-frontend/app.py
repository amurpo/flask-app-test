from flask import Flask, render_template
from prometheus_client import Counter, generate_latest, CONTENT_TYPE_LATEST
from prometheus_client.exposition import choose_encoder

app = Flask(__name__)

# Define a counter metric for HTTP requests
http_requests_total = Counter(
    'http_requests_total', 'Total number of HTTP requests'
)

def create_app():
    # ...
    return app

@app.route('/')
def hello():
    img_url = 'https://oso-flask.s3.sa-east-1.amazonaws.com/oso.jpg'
    message = "Hola, Soy un oso!"

    # Increment the HTTP request counter
    http_requests_total.inc()

    return render_template('index.html', message=message, img_url=img_url)

@app.route('/metrics')
def metrics():
    # Generate and return Prometheus metrics
    encoder, content_type = choose_encoder(['text/plain; version=0.0.4; charset=utf-8'])
    return generate_latest(http_requests_total), 200, {'Content-Type': content_type}

if __name__ == '__main__':
    app = create_app()
    app.run(debug=True)

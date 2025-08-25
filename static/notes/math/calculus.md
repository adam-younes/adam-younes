# Calculus Fundamentals

## Limits

The limit of a function f(x) as x approaches a is L:

$$\lim_{x \to a} f(x) = L$$

## Derivatives

The derivative of a function represents its rate of change:

$$f'(x) = \lim_{h \to 0} \frac{f(x + h) - f(x)}{h}$$

### Common Derivatives

- $\frac{d}{dx}[x^n] = nx^{n-1}$
- $\frac{d}{dx}[\sin(x)] = \cos(x)$
- $\frac{d}{dx}[\cos(x)] = -\sin(x)$
- $\frac{d}{dx}[e^x] = e^x$
- $\frac{d}{dx}[\ln(x)] = \frac{1}{x}$

## Integrals

The integral represents the area under a curve:

$$\int_a^b f(x) dx = F(b) - F(a)$$

Where F(x) is the antiderivative of f(x).

### Common Integrals

- $\int x^n dx = \frac{x^{n+1}}{n+1} + C$ (for n ≠ -1)
- $\int \sin(x) dx = -\cos(x) + C$
- $\int \cos(x) dx = \sin(x) + C$
- $\int e^x dx = e^x + C$
- $\int \frac{1}{x} dx = \ln|x| + C$

## Chain Rule

For composite functions:

$$\frac{d}{dx}[f(g(x))] = f'(g(x)) \cdot g'(x)$$

## Integration by Parts

$$\int u dv = uv - \int v du$$

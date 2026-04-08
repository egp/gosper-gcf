Square Roots via Newton’s Method 

S. G. Johnson, MIT Course 18.335 

February 4, 2015 

1 Overview 

Numerical methods can be distinguished from other branches of analysis and computer science by three characteristics: 

*•* They work with arbitrary real numbers (and vector spaces/extensions thereof): the desired results are not restricted to integers or exact rationals (although in practice we only ever compute *rational approximations* of irrational results). 

*•* Like in computer science (= math \+ time \= math \+ money), we are concerned not only with existence and correctness of the solutions (as in analysis), but with the time (and other computational resources, e.g. memory) required to compute the result. 

*•* We are also concerned with accuracy of the results, because in practice we only ever have approximate answers: 

– Some algorithms may be intrinsically approximate—like the Newton’s-method example shown below, they converge *towards* the desired result but never *reach* it in a finite number of steps. *How fast they converge* is a key question. 

– *Arithmetic with real numbers is approximate* on a computer, because we approximate the set R of real numbers by the set F of floating-point numbers, and the result of every elementary operation (+,*−*,*×*,*÷*) is rounded to the nearest element of F. We need to understand F and how *accumulation* of these rounding errors affects different algorithms. 

2 Square roots 

A classic algorithm that illustrates many of these concerns is “Newton’s” method to compute square roots *x* \=*√~~a~~* for *a \>* 0, i.e. to solve *x*2 \= *a*. The algorithm starts with some guess *x*1 *\>* 0 and computes the sequence of improved guesses   
*xn*\+1 \=12 *xn* \+*axn* *.*   
The intuition is very simple: if *xn* is too big (*\>√~~a~~*), then *a/xn* will be too small (*\<√~~a~~*), and   
so their arithmetic mean *xn*\+1 will be closer to *√~~a~~*. It turns out that this algorithm is very old, dating at least to the ancient Babylonians circa 1000 BCE.1In modern times, this was seen to 

1See e.g. Boyer, *A History of Mathematics*, ch. 3; the Babylonians used base 60 and a famous tablet (YBC 7289\) shows *√*~~2~~ to about six decimal digits. 

1  
be equivalent to Newton’s method to find a root of *f*(*x*) \= *x*2 *− a*. Recall that Newton’s method finds an approximate root of *f*(*x*) \= 0 from a guess *xn* by approximating *f*(*x*) as its tangent line *f*(*xn*) \+ *f0*(*xn*) (*x − xn*), leading to an improved guess *xn*\+1 from the root of the tangent: 

*xn*\+1 \= *xn −f*(*xn*)   
*f0*(*xn*)*,* 

and for *f*(*x*) \= *x*2 *− a* this yields the Babylonian formula above. 

2.1 Convergence proof 

A classic analysis text (Rudin, *Principles of Mathematical Analysis*) approaches the proof of con vergence of this algorithm as follows: we prove that the sequence converges monotonically and is bounded, and hence it has a limit; we then easily see that the limit is *√~~a~~*. In particular: 1\. Suppose *xn \>√~~a~~*, then it follows *√~~a~~ \< xn*\+1 *\< xn*: 

(a) *xn*\+1 *− xn* \=12(*axn− xn*) \= *a−x*2*n*   
2*xn\<* 0\.   
(b) *x*2*n*\+1 *− a* \=14(*x*2*n* \+ 2*a* \+*a*2 

*x*~~2~~*n*) \= 14(*xn −axn*)2 \=(*x*2*n−a*)2   
*x*~~2~~*n*) *− a* \=14(*x*2*n −* 2*a* \+*a*2   
(regardless of whether *xn \>√~~a~~*). 

4*x*~~2~~*n\>* 0 

2\. A monotonic-decreasing sequence that is bounded below converges (Rudin theorem 3.14). If *x*1 *\<√~~a~~*, the second property above means that *x*2 *\>√~~a~~*; then for *n \>* 2 it is monotonically decreasing and bounded below by *√~~a~~*. 

3\. The limit *x* \= lim*n→∞ xn* satisfies *x* \=12(*x* \+*ax*), which is easily solved to show that *x*2 \= *a*. However, this proof by itself tells us nothing about *how fast* the sequence converges 

2.2 Convergence example 

Using the accompanying Julia notebook, we will apply this method to compute the most famous root of all, *√*2\. (Supposedly, the Greek who discovered that *√*2 is irrational was thrown off a cliff by his Pythagorean colleagues.). As a starting guess, we will use *x*1 \= 1, producing the following sequence when computed with about 60 digits of accuracy, where the correct digits are shown in boldface:   
**1** 

**1**.5 

**1**.**41**66666666666666666666666666666666666666666666666666666666675 **1.41421**56862745098039215686274509803921568627450980392156862745 **1.41421356237**46899106262955788901349101165596221157440445849057 **1.41421356237309504880168**96235025302436149819257761974284982890 **1.41421356237309504880168872420969807856967187537**72340015610125 **1.4142135623730950488016887242096980785696718753769480731766796** 

Looking carefully, we see that the number of accurate digits approximately doubles on each iteration. This fantastic convergence rate means that we only need seven Newton iterations to obtain more than 60 accurate digits—the accuracy is quickly limited only by the precision of our floating-point numbers, a topic we will discuss in more detail later on. 

2  
2.3 Convergence rate 

Let us analyze the convergence rate quantitatively—given a small error *δn* on the *n*\-th iteration, we will determine how much smaller the error *δn*\+1 is in the next iteration.   
In particular, let us define *xn* \= *x* (1 \+*δn*), where *x* \=*√~~a~~* is the exact solution. This corresponds to defining *|δn|* as the relative error: 

*|δn|* \=*|xn − x|*   
*|x|,* 

also called the fractional error (the error as a fraction of the exact value). Relative error is typically the most useful way to quantify the error because it is a *dimensionless* quantity (independent of the units or overal scalling of *x*). The logarithm (*−* log10 *δn*) of the relative error is roughly the number of accurate significant digits in the answer *xn*.   
We can plug this definition of *xn* (and *xn*\+1) in terms of *δn* (and *δn*\+1) into our Newton iteration formula to solve for the iteration of *δn*, using the fact that *a/x* \= *x* to divide both sides by *x*: 

1 \+ *δn*\+1 \=12 1 \+ *δn* \+1 


1 \+ *δn* 

\=12 1 \+ *δn* \+ 1 *− δn* \+ *δ*2*n* \+ *O*(*δ*3*n*) *,* 

where we have Taylor-expanded (1*−δn*)*−*1. The *O*(*δ*3*n*) means roughly “terms of order *δ*3*n* or smaller;” we will define it more precisely later on. Because the sequence converges, we are entitled to assume that *|δn|*3   1 for sufficiently large *n*, and so the *δ*3*n* and higher-order terms are eventually negligible   
compared to *δ*2*n*. We obtain: 

*δn*\+1 \=*δ*2*n*2\+ *O*(*δ*3*n*)*,* 

which means the error roughly squares (and halves) on each iteration once we are close to the solution. Squaring the relative error corresponds precisely to doubling the number of significant digits, and hence explains the phenomenon above. This is known as quadratic convergence (not to be confused with “second-order” convergence, which unfortunately refers to an entirely different concept). 


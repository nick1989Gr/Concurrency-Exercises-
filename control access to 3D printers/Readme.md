# Exercise description 
During a busy hackaton 3D printers are heavily used. Write a simulation of hackers trying to access 3D printers.
In the hackerspace, there are 3 3D printers. There are 7 hackers that are interested in using the printers.
If the hacker can't access the printer for more than 5 seconds,
he gets annoyed and quits the hackaton. Hackers use printers for random interval
from 1 to 10 seconds and usually they need to use the printer at least twice, because nothing is perfect for the first time.

Note: This exercise corresponds to exercise 5 from [here](https://tinystruggles.com/2015/10/21/golang-concurrency.html)

# Regarding the solution 

We don't have race conditions because each hacker thread is operating on one element of the hackers array 
Printers thread are in sync with the hacker threads so for a specific hacker struct a printer cannot write at 
the same time that the hacker thread is reading or vice versa 

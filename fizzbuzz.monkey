let fb = fn(n) {
  if (n > 1) {
    fb(n - 1);
  };

  if (n % 15 == 0) {
    puts("FizzBuzz");
    return;
  }
  if (n % 3 == 0) {
    puts("Fuzz");
    return;
  }
  if (n % 5 == 0) {
    puts("Buzz");
    return;
  }
  puts(n);
};

fb(100);

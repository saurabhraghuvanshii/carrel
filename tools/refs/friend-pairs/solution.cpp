int pairings(int n) {
    const long long mod = 1000000007;
    long long a = 1, b = 1;
    for (long long i = 2; i <= n; i++) {
        long long c = (b + (i - 1) * a) % mod;
        a = b;
        b = c;
    }
    return (int)b;
}

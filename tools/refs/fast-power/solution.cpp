int fastPower(long long a, long long e) {
    const long long mod = 1000000007;
    long long result = 1;
    a %= mod;
    for (; e > 0; e >>= 1) {
        if (e & 1) result = result * a % mod;
        a = a * a % mod;
    }
    return result;
}

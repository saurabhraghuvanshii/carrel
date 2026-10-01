// Ported from the owner's bit manipulation/Oddoreven.java (fastExpo), with
// long parameters so the inputs fit. Fixed: there was no modulo, so the
// powers overflowed long; every product now takes the modulo, and the base
// is reduced first.
class Solution {
    public int fastPower(long a, long n) {
        final long mod = 1_000_000_007;
        long ans = 1;
        a %= mod;
        while (n > 0) {
            if ((n & 1) != 0) {
                ans = ans * a % mod;
            }
            a = a * a % mod;
            n = n >> 1;
        }
        return (int) ans;
    }
}

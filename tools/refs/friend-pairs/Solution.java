// Ported from the owner's recursion/FriendPairing.java. Fixed: the recursion
// had no modulo, so the count overflowed an int from n = 19, and it called itself
// twice per level, which takes exponential time. The same recurrence now runs
// bottom-up with a long and the modulo.
class Solution {
    public int pairings(int n) {
        final long mod = 1_000_000_007;
        long prev = 1, cur = 1; // the counts for 0 and 1 people
        for (int i = 2; i <= n; i++) {
            long next = (cur + (i - 1) * prev) % mod;
            prev = cur;
            cur = next;
        }
        return (int) cur;
    }
}

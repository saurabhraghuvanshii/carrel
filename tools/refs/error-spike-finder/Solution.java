class Solution {
    public int worstBurst(int[] times, int window) {
        int best = 0, lo = 0;
        for (int hi = 0; hi < times.length; hi++) {
            while (times[lo] <= times[hi] - window) {
                lo++;
            }
            best = Math.max(best, hi - lo + 1);
        }
        return best;
    }
}

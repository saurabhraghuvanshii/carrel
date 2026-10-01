class Solution {
    public int squareRoot(int x) {
        long lo = 0;
        long hi = x;
        while (lo < hi) {
            long mid = (lo + hi + 1) / 2;
            if (mid * mid <= x) {
                lo = mid;
            } else {
                hi = mid - 1;
            }
        }
        return (int) lo;
    }
}

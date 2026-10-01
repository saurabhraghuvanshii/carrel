class Solution {
    public int shortestWalk(int[] stops, int start) {
        int lo = start, hi = start;
        for (int s : stops) {
            lo = Math.min(lo, s);
            hi = Math.max(hi, s);
        }
        int left = start - lo, right = hi - start;
        return left + right + Math.min(left, right);
    }
}

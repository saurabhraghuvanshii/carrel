class Solution {
    public double median(int[] first, int[] second) {
        if (first.length > second.length) {
            return median(second, first);
        }
        int m = first.length;
        int n = second.length;
        int half = (m + n + 1) / 2;
        int lo = 0;
        int hi = m;
        while (true) {
            int i = (lo + hi) / 2;
            int j = half - i;
            int leftA = i > 0 ? first[i - 1] : Integer.MIN_VALUE;
            int rightA = i < m ? first[i] : Integer.MAX_VALUE;
            int leftB = j > 0 ? second[j - 1] : Integer.MIN_VALUE;
            int rightB = j < n ? second[j] : Integer.MAX_VALUE;
            if (leftA <= rightB && leftB <= rightA) {
                int leftMax = Math.max(leftA, leftB);
                if ((m + n) % 2 == 1) {
                    return leftMax;
                }
                return (leftMax + (long) Math.min(rightA, rightB)) / 2.0;
            }
            if (leftA > rightB) {
                hi = i - 1;
            } else {
                lo = i + 1;
            }
        }
    }
}

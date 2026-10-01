class Solution {
    public int[] surgeMultipliers(int[] requests, int[] drivers, int k) {
        int[] out = new int[requests.length - k + 1];
        int r = 0, d = 0;
        for (int i = 0; i < requests.length; i++) {
            r += requests[i];
            d += drivers[i];
            if (i >= k) {
                r -= requests[i - k];
                d -= drivers[i - k];
            }
            if (i >= k - 1) {
                out[i - k + 1] = r <= d ? 1 : d == 0 ? 5 : Math.min(5, (r + d - 1) / d);
            }
        }
        return out;
    }
}

class Solution {
    public int shortestAtLeast(int[] nums, int target) {
        int best = 0;
        int sum = 0;
        int start = 0;
        for (int i = 0; i < nums.length; i++) {
            sum += nums[i];
            while (sum >= target) {
                if (best == 0 || i - start + 1 < best) {
                    best = i - start + 1;
                }
                sum -= nums[start++];
            }
        }
        return best;
    }
}

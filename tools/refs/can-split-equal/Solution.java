// Ported from the owner's dp/knapscak/TargetSum.java (targetSumSubset), asked
// for half the total.
class Solution {
    public boolean canSplit(int[] nums) {
        int total = 0;
        for (int v : nums) {
            total += v;
        }
        return total % 2 == 0 && targetSumSubset(nums, total / 2);
    }

    private boolean targetSumSubset(int arr[], int sum) {
        int n = arr.length;
        boolean dp[][] = new boolean[n + 1][sum + 1];
        for (int i = 0; i < n + 1; i++) {
            dp[i][0] = true;
        }
        for (int i = 1; i < n + 1; i++) {
            for (int j = 1; j < sum + 1; j++) {
                int v = arr[i - 1];
                if (v <= j && dp[i - 1][j - v] == true) {
                    dp[i][j] = true;
                } else if (dp[i - 1][j] == true) {
                    dp[i][j] = true;
                }
            }
        }
        return dp[n][sum];
    }
}

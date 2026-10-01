class Solution {
    public int[] sortedPair(int[] nums, int target) {
        int i = 0;
        int j = nums.length - 1;
        while (i < j) {
            int sum = nums[i] + nums[j];
            if (sum == target) {
                return new int[] {i, j};
            }
            if (sum < target) {
                i++;
            } else {
                j--;
            }
        }
        return new int[0];
    }
}

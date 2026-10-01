class Solution {
    public int removeRepeats(int[] nums) {
        int k = 0;
        for (int v : nums) {
            if (k == 0 || nums[k - 1] != v) {
                nums[k++] = v;
            }
        }
        return k;
    }
}

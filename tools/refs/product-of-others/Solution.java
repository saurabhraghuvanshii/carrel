class Solution {
    public int[] productOfOthers(int[] nums) {
        int[] out = new int[nums.length];
        int left = 1;
        for (int i = 0; i < nums.length; i++) {
            out[i] = left;
            left *= nums[i];
        }
        int right = 1;
        for (int i = nums.length - 1; i >= 0; i--) {
            out[i] *= right;
            right *= nums[i];
        }
        return out;
    }
}

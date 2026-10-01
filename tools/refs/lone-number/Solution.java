class Solution {
    public int loneNumber(int[] nums) {
        int x = 0;
        for (int v : nums) {
            x ^= v;
        }
        return x;
    }
}

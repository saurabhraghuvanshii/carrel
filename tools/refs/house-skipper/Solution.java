class Solution {
    public int maxCollect(int[] nums) {
        int take = 0, skip = 0;
        for (int v : nums) {
            int next = skip + v;
            skip = Math.max(take, skip);
            take = next;
        }
        return Math.max(take, skip);
    }
}

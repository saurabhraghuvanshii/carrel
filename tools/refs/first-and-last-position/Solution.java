class Solution {
    public int[] firstAndLast(int[] nums, int target) {
        int first = edge(nums, target, true);
        if (first < 0) {
            return new int[] {-1, -1};
        }
        return new int[] {first, edge(nums, target, false)};
    }

    private int edge(int[] nums, int target, boolean leftmost) {
        int lo = 0;
        int hi = nums.length - 1;
        int found = -1;
        while (lo <= hi) {
            int mid = lo + (hi - lo) / 2;
            if (nums[mid] == target) {
                found = mid;
                if (leftmost) {
                    hi = mid - 1;
                } else {
                    lo = mid + 1;
                }
            } else if (nums[mid] < target) {
                lo = mid + 1;
            } else {
                hi = mid - 1;
            }
        }
        return found;
    }
}

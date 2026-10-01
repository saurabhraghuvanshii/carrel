import java.util.*;

class Solution {
    public int[] windowMaximums(int[] nums, int k) {
        int[] out = new int[nums.length - k + 1];
        ArrayDeque<Integer> queue = new ArrayDeque<>();
        for (int i = 0; i < nums.length; i++) {
            while (!queue.isEmpty() && nums[queue.peekLast()] <= nums[i]) {
                queue.pollLast();
            }
            queue.addLast(i);
            if (queue.peekFirst() <= i - k) {
                queue.pollFirst();
            }
            if (i >= k - 1) {
                out[i - k + 1] = nums[queue.peekFirst()];
            }
        }
        return out;
    }
}

import java.util.*;

class Solution {
    public int kthLargest(int[] nums, int k) {
        PriorityQueue<Integer> top = new PriorityQueue<>();
        for (int v : nums) {
            top.add(v);
            if (top.size() > k) {
                top.poll();
            }
        }
        return top.peek();
    }
}

class Solution {
    public int trappedWater(int[] heights) {
        int i = 0;
        int j = heights.length - 1;
        int leftMax = 0;
        int rightMax = 0;
        int total = 0;
        while (i <= j) {
            if (leftMax <= rightMax) {
                leftMax = Math.max(leftMax, heights[i]);
                total += leftMax - heights[i++];
            } else {
                rightMax = Math.max(rightMax, heights[j]);
                total += rightMax - heights[j--];
            }
        }
        return total;
    }
}

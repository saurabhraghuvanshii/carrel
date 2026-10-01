class Solution {
    public int mostWater(int[] heights) {
        int best = 0;
        int i = 0;
        int j = heights.length - 1;
        while (i < j) {
            best = Math.max(best, Math.min(heights[i], heights[j]) * (j - i));
            if (heights[i] < heights[j]) {
                i++;
            } else {
                j--;
            }
        }
        return best;
    }
}

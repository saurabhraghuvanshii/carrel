class Solution {
    public int clashes(int[][] meetings, int start, int end) {
        int count = 0;
        for (int[] m : meetings) {
            if (m[0] < end && start < m[1]) {
                count++;
            }
        }
        return count;
    }
}

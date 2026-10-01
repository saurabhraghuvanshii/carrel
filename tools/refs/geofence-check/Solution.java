class Solution {
    public int entries(int[] zone, int[][] path) {
        int count = 0;
        boolean wasInside = false;
        for (int[] p : path) {
            boolean inside = zone[0] <= p[0] && p[0] <= zone[2] && zone[1] <= p[1] && p[1] <= zone[3];
            if (inside && !wasInside) {
                count++;
            }
            wasInside = inside;
        }
        return count;
    }
}

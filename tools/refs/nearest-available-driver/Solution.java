class Solution {
    public int nearestDriver(int[][] drivers, int rx, int ry) {
        int best = -1;
        long bestDist = Long.MAX_VALUE;
        for (int[] d : drivers) {
            if (d[3] == 0) {
                continue;
            }
            long dist = Math.abs((long) d[1] - rx) + Math.abs((long) d[2] - ry);
            if (dist < bestDist || (dist == bestDist && d[0] < best)) {
                best = d[0];
                bestDist = dist;
            }
        }
        return best;
    }
}

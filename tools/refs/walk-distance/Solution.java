// Ported from the owner's string/Shortestdistance.java (getShortestPath),
// returning x² + y² instead of its square root. Fixed: x * x was an int,
// which overflows once x passes 46340 (the original had the same bug before
// its square root); the squares are now longs.
class Solution {
    public long squaredDistance(String path) {
        int x = 0, y = 0;
        for (int i = 0; i < path.length(); i++) {
            char dir = path.charAt(i);
            if (dir == 'S') {
                y--;
            } else if (dir == 'N') {
                y++;
            } else if (dir == 'W') {
                x--;
            } else {
                x++;
            }
        }
        long X2 = (long) x * x;
        long Y2 = (long) y * y;
        return X2 + Y2;
    }
}

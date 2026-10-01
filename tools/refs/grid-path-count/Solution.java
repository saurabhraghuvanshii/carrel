class Solution {
    public int countPaths(int[][] grid) {
        int rows = grid.length, cols = grid[0].length;
        int[][] ways = new int[rows][cols];
        for (int r = 0; r < rows; r++) {
            for (int c = 0; c < cols; c++) {
                if (grid[r][c] == 1) {
                    continue;
                }
                if (r == 0 && c == 0) {
                    ways[r][c] = 1;
                    continue;
                }
                ways[r][c] = (r > 0 ? ways[r - 1][c] : 0) + (c > 0 ? ways[r][c - 1] : 0);
            }
        }
        return ways[rows - 1][cols - 1];
    }
}

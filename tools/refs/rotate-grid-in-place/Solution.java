class Solution {
    public void rotateGrid(int[][] grid) {
        int n = grid.length;
        for (int r = 0; r < n; r++) {
            for (int c = r + 1; c < n; c++) {
                int t = grid[r][c];
                grid[r][c] = grid[c][r];
                grid[c][r] = t;
            }
        }
        for (int[] row : grid) {
            for (int i = 0, j = n - 1; i < j; i++, j--) {
                int t = row[i];
                row[i] = row[j];
                row[j] = t;
            }
        }
    }
}

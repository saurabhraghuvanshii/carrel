import java.util.*;

class Solution {
    public int minutesToInfectAll(int[][] grid) {
        int rows = grid.length;
        int cols = grid[0].length;
        ArrayDeque<int[]> queue = new ArrayDeque<>();
        int healthy = 0;
        for (int r = 0; r < rows; r++) {
            for (int c = 0; c < cols; c++) {
                if (grid[r][c] == 2) {
                    queue.add(new int[] {r, c});
                } else if (grid[r][c] == 1) {
                    healthy++;
                }
            }
        }
        int minutes = 0;
        int[][] steps = {{0, 1}, {1, 0}, {0, -1}, {-1, 0}};
        while (healthy > 0 && !queue.isEmpty()) {
            for (int size = queue.size(); size > 0; size--) {
                int[] p = queue.poll();
                for (int[] d : steps) {
                    int r = p[0] + d[0];
                    int c = p[1] + d[1];
                    if (r >= 0 && r < rows && c >= 0 && c < cols && grid[r][c] == 1) {
                        grid[r][c] = 2;
                        healthy--;
                        queue.add(new int[] {r, c});
                    }
                }
            }
            minutes++;
        }
        return healthy == 0 ? minutes : -1;
    }
}

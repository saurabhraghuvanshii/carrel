import java.util.*;

class Solution {
    public int shortestPath(int[][] grid) {
        int rows = grid.length;
        int cols = grid[0].length;
        if (grid[0][0] != 0 || grid[rows - 1][cols - 1] != 0) {
            return -1;
        }
        int[][] dist = new int[rows][cols];
        dist[0][0] = 1;
        ArrayDeque<int[]> queue = new ArrayDeque<>();
        queue.add(new int[] {0, 0});
        int[][] steps = {{0, 1}, {1, 0}, {0, -1}, {-1, 0}};
        while (!queue.isEmpty()) {
            int[] p = queue.poll();
            for (int[] d : steps) {
                int r = p[0] + d[0];
                int c = p[1] + d[1];
                if (r >= 0 && r < rows && c >= 0 && c < cols && grid[r][c] == 0 && dist[r][c] == 0) {
                    dist[r][c] = dist[p[0]][p[1]] + 1;
                    queue.add(new int[] {r, c});
                }
            }
        }
        return dist[rows - 1][cols - 1] == 0 ? -1 : dist[rows - 1][cols - 1];
    }
}

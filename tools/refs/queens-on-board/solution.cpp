#include <string>
#include <vector>
using namespace std;

static int place(const vector<string>& b, int row, vector<bool>& col, vector<bool>& d1, vector<bool>& d2) {
    int n = b.size();
    if (row == n) return 1;
    int total = 0;
    for (int c = 0; c < n; c++) {
        if (b[row][c] == '#' || col[c] || d1[row + c] || d2[row - c + n]) continue;
        col[c] = d1[row + c] = d2[row - c + n] = true;
        total += place(b, row + 1, col, d1, d2);
        col[c] = d1[row + c] = d2[row - c + n] = false;
    }
    return total;
}

int countPlacements(vector<string>& board) {
    int n = board.size();
    vector<bool> col(n), d1(2 * n), d2(2 * n + 1);
    return place(board, 0, col, d1, d2);
}

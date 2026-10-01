#include <string>
#include <vector>
using namespace std;

static void walk(const string& w, size_t start, vector<string>& cur, vector<vector<string>>& out) {
    if (start == w.size()) {
        out.push_back(cur);
        return;
    }
    for (size_t end = start + 1; end <= w.size(); end++) {
        string piece = w.substr(start, end - start);
        if (piece == string(piece.rbegin(), piece.rend())) {
            cur.push_back(piece);
            walk(w, end, cur, out);
            cur.pop_back();
        }
    }
}

vector<vector<string>> cuts(string& word) {
    vector<vector<string>> out;
    vector<string> cur;
    walk(word, 0, cur, out);
    return out;
}

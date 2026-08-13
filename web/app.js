let checking = false;
let pollTimer = null;


/* =========================
   THEME
   ========================= */

function getPreferredTheme() {
    const savedTheme =
        localStorage.getItem(
            "vpn-checker-theme"
        );

    if (
        savedTheme === "light" ||
        savedTheme === "dark"
    ) {
        return savedTheme;
    }

    return window.matchMedia(
        "(prefers-color-scheme: dark)"
    ).matches
        ? "dark"
        : "light";
}

function applyTheme(theme) {
    document.documentElement.dataset.theme =
        theme;

    const icon =
        document.getElementById(
            "themeIcon"
        );

    const button =
        document.getElementById(
            "themeToggle"
        );

    if (icon) {
        icon.textContent =
            theme === "dark"
                ? "☀️"
                : "🌙";
    }

    if (button) {
        const label =
            theme === "dark"
                ? "Переключить на светлую тему"
                : "Переключить на тёмную тему";

        button.setAttribute(
            "aria-label",
            label
        );

        button.setAttribute(
            "title",
            label
        );
    }
}

function initTheme() {
    applyTheme(
        getPreferredTheme()
    );

    const button =
        document.getElementById(
            "themeToggle"
        );

    if (!button) {
        return;
    }

    button.addEventListener(
        "click",
        () => {
            const currentTheme =
                document.documentElement
                    .dataset
                    .theme;

            const nextTheme =
                currentTheme === "dark"
                    ? "light"
                    : "dark";

            localStorage.setItem(
                "vpn-checker-theme",
                nextTheme
            );

            applyTheme(nextTheme);
        }
    );
}


/* =========================
   HELPERS
   ========================= */

async function getJSON(url) {
    const response = await fetch(url, {
        cache: "no-store",
    });

    if (!response.ok) {
        throw new Error(
            `HTTP ${response.status}`
        );
    }

    return await response.json();
}

function setText(id, value) {
    const element =
        document.getElementById(id);

    if (element) {
        element.textContent = value;
    }
}

function formatDuration(seconds) {
    const value = Number(seconds);

    if (
        !Number.isFinite(value) ||
        value <= 0
    ) {
        return "—";
    }

    const totalSeconds =
        Math.ceil(value);

    const hours =
        Math.floor(
            totalSeconds / 3600
        );

    const minutes =
        Math.floor(
            (totalSeconds % 3600) / 60
        );

    const secs =
        totalSeconds % 60;

    if (hours > 0) {
        return `${hours} ч ${minutes} мин`;
    }

    if (minutes > 0) {
        return `${minutes} мин ${secs} сек`;
    }

    return `${secs} сек`;
}

function formatSpeed(speed) {
    const value = Number(speed);

    if (
        !Number.isFinite(value) ||
        value <= 0
    ) {
        return "—";
    }

    return `${value.toFixed(1)} cfg/s`;
}

function updateCards(
    total,
    working,
    failed
) {
    setText(
        "total",
        total
    );

    setText(
        "working",
        working
    );

    setText(
        "failed",
        failed
    );
}


/* =========================
   STATS
   ========================= */

async function updateStats() {
    try {
        const data =
            await getJSON(
                "/api/stats?t=" +
                Date.now()
            );

        const container =
            document.getElementById(
                "protocols"
            );

        if (!container) {
            return;
        }

        const protocols =
            data.protocols || {};

        container.innerHTML =
            Object.entries(protocols)
                .map(
                    ([name, stat]) => `
                        <p>
                            <strong>${name}</strong>:
                            ${stat.working ?? 0}/${
                                stat.total ?? 0
                            }
                            рабочих
                        </p>
                    `
                )
                .join("");

    } catch (error) {
        console.error(
            "updateStats:",
            error
        );
    }
}


/* =========================
   BEST CONFIGS
   ========================= */

async function updateBest() {
    try {
        const data =
            await getJSON(
                "/api/best?limit=5&t=" +
                Date.now()
            );

        const container =
            document.getElementById(
                "best"
            );

        if (!container) {
            return;
        }

        container.innerHTML =
            data
                .map(
                    (item, index) => `
                        <p>
                            <strong>
                                #${index + 1}
                            </strong>
                            ${item.protocol ?? ""}
                            —
                            ${item.address ?? ""}:${
                                item.port ?? ""
                            }
                            —
                            ${item.latency_ms ?? 0} ms
                        </p>
                    `
                )
                .join("");

    } catch (error) {
        console.error(
            "updateBest:",
            error
        );
    }
}


/* =========================
   POLLING
   ========================= */

function startPolling() {
    if (pollTimer !== null) {
        return;
    }

    pollTimer = setInterval(
        updateStatus,
        1000
    );
}

function stopPolling() {
    if (pollTimer !== null) {
        clearInterval(
            pollTimer
        );

        pollTimer = null;
    }
}


/* =========================
   CHECK STATUS
   ========================= */

async function updateStatus() {
    try {
        const data =
            await getJSON(
                "/api/check/status?t=" +
                Date.now()
            );

        const total =
            Number(data.total) || 0;

        const checked =
            Number(data.checked) || 0;

        const working =
            Number(data.working) || 0;

        const failed =
            Number(data.failed) || 0;

        const speed =
            Number(data.current_speed) || 0;

        const eta =
            Number(
                data.estimated_seconds_left
            ) || 0;

        let progress = 0;

        if (total > 0) {
            progress =
                Math.floor(
                    checked * 100 / total
                );
        }

        if (
            data.status === "completed"
        ) {
            progress = 100;
        }

        setText(
            "progress",
            `${progress}%`
        );

        const progressBar =
            document.getElementById(
                "progressBar"
            );

        if (progressBar) {
            progressBar.style.width =
                `${progress}%`;
        }

        if (
            data.status === "running"
        ) {
            updateCards(
                total,
                working,
                failed
            );

            const speedText =
                formatSpeed(speed);

            const etaText =
                eta > 0
                    ? `осталось ${formatDuration(eta)}`
                    : "расчёт ETA...";

            setText(
                "status",
                `Проверка — ${checked}/${total}` +
                ` • ${speedText}` +
                ` • ${etaText}`
            );

            const button =
                document.getElementById(
                    "checkButton"
                );

            if (button) {
                button.disabled = true;
                button.textContent =
                    "Проверка...";
            }

            checking = true;

            startPolling();

            return data;
        }

        if (
            data.status === "completed"
        ) {
            updateCards(
                total,
                working,
                failed
            );

            setText(
                "progress",
                "100%"
            );

            if (progressBar) {
                progressBar.style.width =
                    "100%";
            }

            setText(
                "status",
                `Завершено — ${checked}/${total}` +
                ` • ${formatSpeed(speed)}`
            );

            const button =
                document.getElementById(
                    "checkButton"
                );

            if (button) {
                button.disabled = false;
                button.textContent =
                    "Запустить проверку";
            }

            checking = false;

            stopPolling();

            await Promise.all([
                updateStats(),
                updateBest(),
            ]);

            return data;
        }

        if (
            data.status === "error"
        ) {
            updateCards(
                total,
                working,
                failed
            );

            setText(
                "status",
                `Ошибка — ${checked}/${total}`
            );

            const button =
                document.getElementById(
                    "checkButton"
                );

            if (button) {
                button.disabled = false;
                button.textContent =
                    "Запустить проверку";
            }

            checking = false;

            stopPolling();

            await Promise.all([
                updateStats(),
                updateBest(),
            ]);

            return data;
        }

        setText(
            "status",
            `${data.status} — ${checked}/${total}`
        );

        stopPolling();

        return data;

    } catch (error) {
        console.error(
            "updateStatus:",
            error
        );

        return null;
    }
}


/* =========================
   START CHECK
   ========================= */

async function startCheck() {
    if (checking) {
        return;
    }

    checking = true;

    const button =
        document.getElementById(
            "checkButton"
        );

    setText(
        "progress",
        "0%"
    );

    setText(
        "status",
        "Запуск проверки..."
    );

    const progressBar =
        document.getElementById(
            "progressBar"
        );

    if (progressBar) {
        progressBar.style.width =
            "0%";
    }

    if (button) {
        button.disabled = true;
        button.textContent =
            "Запуск...";
    }

    try {
        const response =
            await fetch(
                "/api/check",
                {
                    method: "POST",
                    cache: "no-store",
                }
            );

        if (!response.ok) {
            throw new Error(
                `HTTP ${response.status}`
            );
        }

        await updateStatus();

        startPolling();

    } catch (error) {
        console.error(
            "startCheck:",
            error
        );

        checking = false;

        stopPolling();

        if (button) {
            button.disabled = false;
            button.textContent =
                "Запустить проверку";
        }

        setText(
            "status",
            "Ошибка запуска проверки"
        );
    }
}


/* =========================
   INIT
   ========================= */

async function init() {
    initTheme();

    const button =
        document.getElementById(
            "checkButton"
        );

    if (button) {
        button.addEventListener(
            "click",
            startCheck
        );
    }

    await Promise.all([
        updateStatus(),
        updateStats(),
        updateBest(),
    ]);
}

document.addEventListener(
    "DOMContentLoaded",
    init
);